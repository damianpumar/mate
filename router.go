package mate

import (
	"fmt"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type Router struct {
	router *http.ServeMux
}

func NewRouter() *Router {
	return &Router{router: http.NewServeMux()}
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(res, req)
}

func (s *Router) Routes() *http.ServeMux {
	return s.router
}

func (r *Router) Get(path string, handler func(c *Context)) {
	r.addRoute(http.MethodGet, path, handler)
}

func (r *Router) Post(path string, handler func(c *Context)) {
	r.addRoute(http.MethodPost, path, handler)
}

func (r *Router) Put(path string, handler func(c *Context)) {
	r.addRoute(http.MethodPut, path, handler)
}

func (r *Router) Delete(path string, handler func(c *Context)) {
	r.addRoute(http.MethodDelete, path, handler)
}

func (r *Router) Patch(path string, handler func(c *Context)) {
	r.addRoute(http.MethodPatch, path, handler)
}

func (r *Router) Folder(path string, dir string) {
	fileserver := http.FileServer(http.Dir(dir))

	if path != "" {
		fileserver = http.StripPrefix(path, fileserver)
	}

	r.addRoute(http.MethodGet, path, func(c *Context) {
		fileserver.ServeHTTP(c.Response.ResponseWriter, c.Request.Request)
	})
}

func (r *Router) File(path string, file string) {
	r.addRoute(http.MethodGet, path, func(c *Context) {
		http.ServeFile(c.Response.ResponseWriter, c.Request.Request, file)
	})
}

func (r *Router) addRoute(method string, path string, handler func(c *Context)) {
	pattern := fmt.Sprintf("%s %s", method, path)

	if len(path) > 1 && path[len(path)-1] == '/' {
		pattern = strings.TrimSuffix(pattern, "/")
	}

	r.router.HandleFunc(
		pattern, func(w http.ResponseWriter, req *http.Request) {
			handler(NewContext(w, req))
		})
}
