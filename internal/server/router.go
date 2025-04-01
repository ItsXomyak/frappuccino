package server

import (
	"net/http"
)

func newRouter(serv *Handler) *http.ServeMux {
	router := http.NewServeMux()

	RegisterRoutes(router, serv)
	return router
}
