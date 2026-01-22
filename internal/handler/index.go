package handler

import (
	"net/http"

	"github.com/gngtwhh/WBlog/internal/render"
	"github.com/gngtwhh/WBlog/internal/service"
	"github.com/gngtwhh/WBlog/pkg/errcode"
	"github.com/gngtwhh/WBlog/pkg/response"
)

type IndexData struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type IndexHandler struct {
	articleSvc *service.ArticleService
}

func NewIndexHandler(svc *service.ArticleService) *IndexHandler {
	return &IndexHandler{articleSvc: svc}
}

// Index returns a json data
func (h *IndexHandler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		response.Fail(w, errcode.NotFound)
		return
	}
	h.IndexHtml(w, r)
}

func (h *IndexHandler) IndexHtml(w http.ResponseWriter, r *http.Request) {
	render.Execute(w, "index", nil)
}

func (h *IndexHandler) ArticlePage(w http.ResponseWriter, r *http.Request) {
	render.Execute(w, "article", nil)
}
