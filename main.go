package main

import (
	"flag"
	"minimal/database"
	"minimal/framework"
	"minimal/server"
	"net/http"
	"time"
)

type Example struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

var (
	port = flag.String("port", "8080", "Port to listen on")
)

func main() {
	flag.Parse()

	server := server.New()

	db := database.Connect()

	cookie := framework.NewSecureCookie("my-secret")

	server.Get("/", framework.LoggingMiddleware(func(c *framework.Context) {
		data := db.Select("users")

		c.JSON(200, data)
	}))

	server.Get("/cookie", framework.LoggingMiddleware(func(c *framework.Context) {
		err := cookie.SetEncryptedCookie(c.Response, "session", "user12345", 30*time.Second)
		if err != nil {
			http.Error(c.Response, err.Error(), http.StatusInternalServerError)
			return
		}

		c.Text(200, "Cookie set")
	}))

	server.Get("/read-cookie", framework.LoggingMiddleware(func(c *framework.Context) {
		value, err := cookie.GetEncryptedCookie(c.Request.Request, "session")
		if err != nil {
			http.Error(c.Response, err.Error(), http.StatusInternalServerError)
			return
		}

		c.Text(200, value)
	}))

	server.Get("/delete-cookie", framework.LoggingMiddleware(func(c *framework.Context) {
		cookie.ClearCookie(c.Response, "session")
		c.Response.Text(200, "Cookie deleted")
	}))

	server.Post("/post", func(c *framework.Context) {
		data := Example{}

		c.BindBody(&data)

		db.Insert("users", data)

		c.JSON(200, data)
	})

	server.Get("/{id}", func(c *framework.Context) {
		id := c.GetPathValue("id")

		data := db.SelectById("users", id)

		c.JSON(200, data)
	})

	server.Start(port)
}
