package framework

import (
	"fmt"
	"net/http"
	"os"
)

type Server struct {
	router *Router
}

func NewServer() *Server {
	return &Server{
		router: NewRouter(),
	}
}

func (s *Server) Start(port *string) {
	routes := s.router.Routes()

	fmt.Println("🚀 Server running on", "http://localhost:"+*port)

	if err := http.ListenAndServe(":"+*port, routes); err != nil {
		fmt.Println("🤔 Error starting server", err)

		os.Exit(1)
	}
}

func (s *Server) Get(path string, handler HandlerFunc) {
	s.router.Get(path, handler)
}
