package framework

import (
	"net/http"
)

type HandlerFunc func(c *Context)

type Router struct {
	router *http.ServeMux
}

func NewRouter() *Router {
	return &Router{router: http.NewServeMux()}
}

func (s *Router) Routes() *http.ServeMux {
	return s.router
}

func (s *Router) Get(path string, handler func(c *Context)) {
	s.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		handler(NewContext(w, r))
	})
}
