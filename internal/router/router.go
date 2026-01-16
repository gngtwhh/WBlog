package router

import (
	"log/slog"
	"net/http"

	"github.com/gngtwhh/WBlog/internal/handler"
	"github.com/gngtwhh/WBlog/internal/middleware"
)

func LoadRouters(app *handler.App, logger *slog.Logger) http.Handler {
	router := http.NewServeMux()

	// static resources
	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../web/static"))))

	// root and /index
	// TODO: use app
	router.HandleFunc("GET /", app.Index.Index)
	router.HandleFunc("GET /index", app.Index.IndexHtml)

	// /admin
	router.HandleFunc("GET /admin", app.Index.Admin)

	// article page
	router.HandleFunc("GET /article/{id}", app.Index.ArticlePage)

	// article api
	router.HandleFunc("GET /api/list-articles", app.Article.ListArticles)
	router.HandleFunc("GET /api/articles-count", app.Article.Count)
	router.HandleFunc("GET /api/get-article", app.Article.GetArticle)

	router.HandleFunc("POST /api/create-article", app.Article.Create)
	router.HandleFunc("POST /api/update-article", app.Article.Update)
	router.HandleFunc("DELETE /api/delete-article", app.Article.Delete)

	// user api
	router.HandleFunc("GET /api/userinfo", app.User.GetUserInfo)
	router.HandleFunc("POST /api/user/register", app.User.Register)
	router.HandleFunc("POST /api/user/login", app.User.Login)
	// authentication required
	{
		router.HandleFunc("GET /api/user/profile", middleware.Auth(app.User.GetProfile))
		router.HandleFunc("POST /api/user/update", middleware.Auth(app.User.UpdateProfile))
		router.HandleFunc("POST /api/user/update-password", middleware.Auth(app.User.UpdatePassword))
		router.HandleFunc("POST /api/user/upload-avatar", middleware.Auth(app.User.UploadAvatar))
	}

	var handler http.Handler = router
	handler = middleware.RequestLogger(logger)(handler)

	return handler
}
