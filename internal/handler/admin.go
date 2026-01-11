package handler

import (
	"net/http"

	"github.com/gngtwhh/WBlog/internal/render"
)

func (h *IndexHandler) Admin(w http.ResponseWriter, r *http.Request) {
	// Render admin.html, no data needed for now,
	// data will be loaded asynchronously via JS
	render.Execute(w, "admin", nil)
}
