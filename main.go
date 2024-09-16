package main

import (
	"flag"
	"minimal/framework"
)

type Example struct {
	Name string
}

var (
	port = flag.String("port", "8080", "Port to listen on")
)

func main() {
	flag.Parse()

	server := framework.NewServer()

	server.Get("/", framework.LoggingMiddleware(func(c *framework.Context) {
		c.Response.WithJson(200, Example{Name: "World"})
	}))

	server.Start(port)
}
