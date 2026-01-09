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
	router.HandleFunc("GET /api/articles-count", app.Article.Count)
	router.HandleFunc("GET /api/get-article", app.Article.GetArticle)
	router.HandleFunc("GET /api/create-article", app.Article.Create)
	router.HandleFunc("GET /api/update-article", app.Article.Update)
	router.HandleFunc("GET /api/delete-article", app.Article.Delete)

	return router
}
