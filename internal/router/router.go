package router

import "net/http"

func LoadRouters() (router *http.ServeMux) {
	router = http.NewServeMux()
	RegisterIndexRoutes(router)
	return router
}
