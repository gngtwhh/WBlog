package router

import (
	"net/http"

	"github.com/gngtwhh/WBlog/internal/handler"
)

func RegisterIndexRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", handler.Index)
	mux.HandleFunc("GET /index", handler.IndexHtml)
}
