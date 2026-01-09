package handler

import (
	"net/http"

	"github.com/gngtwhh/WBlog/internal/service"
)

type ArticleHandler struct {
	svc *service.ArticleService
}

func NewArticleHandler(svc *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{svc: svc}
}

func (h *ArticleHandler) ListArticles(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	// // get param:
	// // @count: count of articles per-page
	// // @page: page index(start from 1)
	// count := r.URL.Query().Get("count")
	// page := r.URL.Query().Get("page")

	// json.NewEncoder(w).Encode(articles)
}
