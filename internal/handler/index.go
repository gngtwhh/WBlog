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

// Article contains information about a blog article.
// TODO: This struct will be moved to model package.
type Article struct {
	Title    string
	Author   string
	Abstract string
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
	// render.Execute(w, "index", nil)
	// TODO: test data to be removed
	articles := []Article{
		{
			Title:    "Go 语言原本",
			Author:   "欧长坤",
			Abstract: "本书讨论 Go 语言核心实现的相关话题...",
		},
		{
			Title:    "Gin 框架入门实战",
			Author:   "WBlog Team",
			Abstract: "Gin 是一个用 Go (Golang) 编写的 HTTP Web 框架...",
		},
		{
			Title:    "为什么选择原生开发",
			Author:   "Gopher",
			Abstract: "原生开发能让你更深刻地理解 HTTP 协议与 Web 原理...",
		},
	}
	testData := map[string]interface{}{
		"Posts": articles,
		"Total": len(articles),
	}
	render.Execute(w, "index", testData)
}
