package handler

import (
	"net/http"

	"github.com/gngtwhh/WBlog/internal/render"
	"github.com/gngtwhh/WBlog/internal/service"
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
	h.IndexHtml(w, r)
	// w.Header().Set("Content-Type", "application/json")
	// // TODO: configurable data
	// data := IndexData{
	// 	Title: "WBlog - A Simple Blog",
	// 	Desc:  "written in golang",
	// }
	// jsonStr, _ := json.Marshal(data)
	// w.Write(jsonStr)
}

func (h *IndexHandler) IndexHtml(w http.ResponseWriter, r *http.Request) {
	render.Execute(w, "index", nil)
}

func (h *IndexHandler) ArticlePage(w http.ResponseWriter, r *http.Request) {
	render.Execute(w, "article", nil)
}
