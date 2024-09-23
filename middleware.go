package mate

type Middleware func(HandlerFunc) HandlerFunc

func LoggingMiddleware(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		println("Request:", c.Request.Method, c.Request.URL.Path)

		next(c)
	}
}
