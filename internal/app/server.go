package app

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gngtwhh/WBlog/internal/render"
	"github.com/gngtwhh/WBlog/internal/router"
)

type Server struct {
	server http.Server
}

func NewServer() (h *Server) {
	// init render
	tmpls := loadTmlps()
	render.Init(tmpls)

	h = &Server{
		server: http.Server{
			Addr: ":8080",
		},
	}
	h.server.Handler = router.LoadRouters()
	return
}

func (s *Server) Run() {
	if err := s.server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func loadTmlps() map[string]*template.Template {
	tmpls := make(map[string]*template.Template)
	base := "../web/templates/"
	layout := base + "layout.html"
	tmpls["index"] = template.Must(template.ParseFiles(layout, base+"index.html"))
	// tmpls["layout"] = template.Must(template.ParseFiles("web/templates/layout.html"))
	return tmpls
}
