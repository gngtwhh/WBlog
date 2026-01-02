package app

import (
	"log"
	"net/http"

	"github.com/gngtwhh/WBlog/internal/router"
)

type Server struct {
	server http.Server
}

func NewServer() (h *Server) {
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
