package framework

import (
	"fmt"
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

func (r *Router) Get(path string, handler func(c *Context)) {
	r.addRoute(http.MethodGet, path, handler)
}

func (s *Router) Post(path string, handler func(c *Context)) {
	s.addRoute(http.MethodPost, path, handler)
}

func (s *Router) Put(path string, handler func(c *Context)) {
	s.addRoute(http.MethodPut, path, handler)
}

func (s *Router) Delete(path string, handler func(c *Context)) {
	s.addRoute(http.MethodDelete, path, handler)
}

func (s *Router) Patch(path string, handler func(c *Context)) {
	s.addRoute(http.MethodPatch, path, handler)
}

func (s *Router) addRoute(method string, path string, handler func(c *Context)) {
	pattern := fmt.Sprintf("%s %s", method, path)

	s.router.HandleFunc(
		pattern, func(w http.ResponseWriter, req *http.Request) {
			handler(NewContext(w, req))
		})
}

func (s *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
