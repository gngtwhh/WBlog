package router

import (
	"net/http"

	"github.com/gngtwhh/WBlog/internal/handler"
)

func LoadRouters(app *handler.App) (router *http.ServeMux) {
	router = http.NewServeMux()

	// root and /index
	// TODO: use app
	router.HandleFunc("GET /", app.Index.Index)
	router.HandleFunc("GET /index", app.Index.IndexHtml)

	// article api
	router.HandleFunc("GET /api/list-articles", app.Article.ListArticles)

	return router
}
