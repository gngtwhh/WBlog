package app

import (
	"bufio"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gngtwhh/WBlog/internal/cache"
	"github.com/gngtwhh/WBlog/internal/config"
	"github.com/gngtwhh/WBlog/internal/handler"
	"github.com/gngtwhh/WBlog/internal/render"
	"github.com/gngtwhh/WBlog/internal/repository"
	"github.com/gngtwhh/WBlog/internal/router"
	"github.com/gngtwhh/WBlog/internal/service"
	"github.com/gngtwhh/WBlog/pkg/logger"
	"github.com/gngtwhh/WBlog/pkg/sensitive"
	"github.com/gngtwhh/WBlog/pkg/utils"
)

type Server struct {
	server http.Server
	logger *slog.Logger
}

func NewServer() (h *Server) {
	// log setup
	log := logger.Setup(&logger.Options{
		Level:     slog.LevelDebug,
		FilePath:  config.Cfg.App.LogFile,
		AddSource: false,
	})
	log.Info("starting WBLOG server",
		slog.Group("config",
			slog.String("mode", config.Cfg.Server.RunMode),
			slog.String("port", config.Cfg.Server.Port),
			slog.String("db_dsn", config.Cfg.Database.DSN),
			slog.String("log_file", config.Cfg.App.LogFile),
		),
	)

	// utils
	// jwt
	if err := utils.InitJwt(config.Cfg.App.JwtSecret); err != nil {
		log.Error("failed to init jwt pkg", "err", err)
		panic(err)
	}
	// sensitive words filter
	file, err := os.Open(config.Cfg.App.SensitiveWordsFile)
	if err != nil {
		log.Error("failed to load sensitive words file", "err", err)
		panic(err)
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words = append(words, word)
		}
	}
	acFilter := sensitive.NewACFilter()
	acFilter.Build(words)

	// init redis cache
	if err := cache.InitRedis(config.Cfg.Cache.RedisAddr, config.Cfg.Cache.RedisPassword); err != nil {
		log.Error("init cache failed", "err", err)
		panic(err)
	}

	// init repository
	log.Info("initializing database...")
	db, err := repository.InitDB(config.Cfg.Database.DSN)
	if err != nil {
		log.Error("failed to connect database", "err", err)
		panic(err)
	}
	articleRepo := repository.NewArticleRepo(db, log)
	userRepo := repository.NewUserRepo(db, log)
	commentRepo := repository.NewCommentRepo(db, log)

	log.Info("initializing service...")
	// init Services
	articleService := service.NewArticleService(articleRepo, log)
	userService := service.NewUserService(userRepo, log)
	commentService := service.NewCommentService(commentRepo, acFilter, log)

	// init handler
	app := &handler.App{
		Index:   handler.NewIndexHandler(articleService),
		Article: handler.NewArticleHandler(articleService),
		User:    handler.NewUserHandler(userService),
		Comment: handler.NewCommentHandler(commentService, articleService),
	}

	// html template pre-compile
	log.Info("pre-compiling html templates...")
	tmpls := loadTmlps()
	render.Init(tmpls, "layout")

	h = &Server{
		server: http.Server{
			Addr:    ":" + config.Cfg.Server.Port,
			Handler: router.LoadRouters(app, log),
		},
		logger: log,
	}
	return
}

func (s *Server) Run() {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error("server startup failed", "err", err)
	}
}

func loadTmlps() map[string]*template.Template {
	tmpls := make(map[string]*template.Template)

	baseDir := config.Cfg.App.TemplateDir
	layout := baseDir + "layout/layout.html"

	tmpls["index"] = template.Must(template.ParseFiles(layout, baseDir+"index.html"))
	tmpls["admin"] = template.Must(template.ParseFiles(layout, baseDir+"admin.html"))
	tmpls["article"] = template.Must(template.ParseFiles(layout, baseDir+"article.html"))
	// tmpls["layout"] = template.Must(template.ParseFiles("web/templates/layout.html"))
	return tmpls
}
