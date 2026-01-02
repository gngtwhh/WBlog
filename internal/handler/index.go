package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gngtwhh/WBlog/internal/render"
)

type IndexData struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

// Index returns a json data
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := IndexData{
		Title: "WBlog - A Simple Blog",
		Desc:  "written in golang",
	}
	jsonStr, _ := json.Marshal(data)
	w.Write(jsonStr)
}

func IndexHtml(w http.ResponseWriter, r *http.Request) {
	render.Execute(w, "index", nil)
}
