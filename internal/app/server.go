package app

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gngtwhh/WBlog/internal/handler"
	"github.com/gngtwhh/WBlog/internal/render"
	"github.com/gngtwhh/WBlog/internal/repository"
	"github.com/gngtwhh/WBlog/internal/router"
	"github.com/gngtwhh/WBlog/internal/service"
)

type Server struct {
	server http.Server
}

func NewServer() (h *Server) {
	tmpls := loadTmlps()
	render.Init(tmpls, "layout")

	// init repository
	db, err := repository.InitDB("./blog.db")
	if err != nil {
		panic(err)
	}
	articleRepo := repository.NewArticleRepo(db)
	userRepo := repository.NewUserRepo(db)

	// init Services
	articleService := service.NewArticleService(articleRepo)
	userService := service.NewUserService(userRepo)

	// init handler
	app := &handler.App{
		Index:   handler.NewIndexHandler(articleService),
		Article: handler.NewArticleHandler(articleService),
		User:    handler.NewUserHandler(userService),
	}

	h = &Server{
		server: http.Server{
			Addr: ":8080",
		},
	}
	h.server.Handler = router.LoadRouters(app)
	return
}

func (s *Server) Run() {
	if err := s.server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func loadTmlps() map[string]*template.Template {
	tmpls := make(map[string]*template.Template)

	baseDir := "../web/templates/"
	layout := baseDir + "layout/layout.html"

	tmpls["index"] = template.Must(template.ParseFiles(layout, baseDir+"index.html"))
	tmpls["admin"] = template.Must(template.ParseFiles(layout, baseDir+"admin.html"))
	tmpls["article"] = template.Must(template.ParseFiles(layout, baseDir+"article.html"))
	// tmpls["layout"] = template.Must(template.ParseFiles("web/templates/layout.html"))
	return tmpls
}
