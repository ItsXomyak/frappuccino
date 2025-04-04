package server

import (
	"log"
	"net/http"
)

type Server struct {
	addr   string
	router *http.ServeMux
}

func enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("CORS middleware: %s %s", r.Method, r.URL.Path)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			log.Println("Handling OPTIONS preflight request")
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func NewServer(port string, serv Handler) *Server {
	router := http.NewServeMux()
	RegisterRoutes(router, &serv)

	addr := ":" + port
	return &Server{
		addr:   addr,
		router: router,
	}
}

func (s *Server) Start() {
	log.Printf("Starting server on %s...\n", s.addr)
	handlerWithCORS := enableCORS(s.router)
	if err := http.ListenAndServe(s.addr, handlerWithCORS); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
