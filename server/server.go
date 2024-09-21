package server

import (
	"fmt"
	"minimal/framework"
	"net/http"
	"os"
)

type Server struct {
	router *framework.Router
}

func New() *Server {
	return &Server{
		router: framework.NewRouter(),
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

func (s *Server) Get(path string, handler framework.HandlerFunc) {
	s.router.Get(path, handler)
}

func (s *Server) Post(path string, handler framework.HandlerFunc) {
	s.router.Post(path, handler)
}

func (s *Server) Put(path string, handler framework.HandlerFunc) {
	s.router.Put(path, handler)
}

func (s *Server) Delete(path string, handler framework.HandlerFunc) {
	s.router.Delete(path, handler)
}
