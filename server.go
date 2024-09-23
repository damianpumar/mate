package mate

import (
	"fmt"
	"net/http"
	"os"
)

type Server struct {
	router *Router
}

func New() *Server {
	return &Server{
		router: NewRouter(),
	}
}

func (s *Server) Start(port string) {
	routes := s.router.Routes()

	fmt.Println("🧉 Server running on", "http://localhost:"+port)

	if err := http.ListenAndServe(":"+port, routes); err != nil {
		fmt.Println("🤔 Error starting server", err)

		os.Exit(1)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Get(path string, handler HandlerFunc) {
	s.router.Get(path, handler)
}

func (s *Server) Post(path string, handler HandlerFunc) {
	s.router.Post(path, handler)
}

func (s *Server) Put(path string, handler HandlerFunc) {
	s.router.Put(path, handler)
}

func (s *Server) Delete(path string, handler HandlerFunc) {
	s.router.Delete(path, handler)
}

func (s *Server) Patch(path string, handler HandlerFunc) {
	s.router.Patch(path, handler)
}

type Group struct {
	path   string
	router *Router
}

func (s *Server) Group(path string, group func(r *Group)) {
	group(&Group{path, s.router})
}

func (g *Group) Get(path string, handler HandlerFunc) {
	groupedPath := g.path + path
	g.router.Get(groupedPath, handler)
}

func (g *Group) Post(path string, handler HandlerFunc) {
	groupedPath := g.path + path
	g.router.Post(groupedPath, handler)
}

func (g *Group) Put(path string, handler HandlerFunc) {
	groupedPath := g.path + path
	g.router.Put(groupedPath, handler)
}

func (g *Group) Delete(path string, handler HandlerFunc) {
	groupedPath := g.path + path
	g.router.Delete(groupedPath, handler)
}

func (g *Group) Patch(path string, handler HandlerFunc) {
	groupedPath := g.path + path
	g.router.Patch(groupedPath, handler)
}
