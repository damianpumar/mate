package main

import (
	"flag"
	"minimal/framework"
	"net/http"
	"time"
)

type Example struct {
	Name string
}

var (
	port = flag.String("port", "8080", "Port to listen on")
)

func main() {
	flag.Parse()

	cookie := framework.NewSecureCookie("my-secret")
	server := framework.NewServer()

	server.Get("/", framework.LoggingMiddleware(func(c *framework.Context) {
		c.Response.WithJson(200, Example{Name: "World"})
	}))

	server.Get("/cookie", framework.LoggingMiddleware(func(c *framework.Context) {
		err := cookie.SetEncryptedCookie(c.Response, "session", "user12345", 30*time.Second)
		if err != nil {
			http.Error(c.Response, err.Error(), http.StatusInternalServerError)
			return
		}

		c.Response.WithText(200, "Cookie set")
	}))

	server.Get("/read-cookie", framework.LoggingMiddleware(func(c *framework.Context) {
		value, err := cookie.GetEncryptedCookie(c.Request.Request, "session")
		if err != nil {
			http.Error(c.Response, err.Error(), http.StatusInternalServerError)
			return
		}

		c.Response.WithText(200, value)
	}))

	server.Get("/delete-cookie", framework.LoggingMiddleware(func(c *framework.Context) {
		cookie.ClearCookie(c.Response, "session")
		c.Response.WithText(200, "Cookie deleted")
	}))

	server.Post("/post", func(c *framework.Context) {
		data := &Example{}

		c.Request.ParseBody(data)

		c.Response.WithJson(200, data)
	})

	server.Start(port)
}
