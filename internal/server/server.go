package server

import (
	"log"
	"net/http"
)

type Server struct {
	addr   string
	router *http.ServeMux
}

func NewServer(port string, serv Handler) *Server {
	router := newRouter(&serv)

	addr := ":" + port
	return &Server{
		addr:   addr,
		router: router,
	}
}

func (s *Server) Start() {
	log.Printf("Starting server on %s...\n", s.addr)
	if err := http.ListenAndServe(s.addr, s.router); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
