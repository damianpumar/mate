package main

import (
	"flag"
	"mate/database"
	"mate/framework"
	"mate/http"
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

	server := http.New()

	db := database.Connect()

	cookie := framework.NewSecureCookie("my-secret")

	server.Get("/cookie", framework.LoggingMiddleware(func(c *framework.Context) {
		if err := cookie.SetEncryptedCookie(c.Response, "session", "user12345", 30*time.Second); err != nil {

			c.Error(500, err)

			return
		}

		c.Text(200, "Cookie set")
	}))

	server.Get("/read-cookie", framework.LoggingMiddleware(func(c *framework.Context) {
		value, err := cookie.GetEncryptedCookie(c.Request.Request, "session")

		if err != nil {
			c.Error(500, err)

			return
		}

		c.Text(200, value)
	}))

	server.Get("/delete-cookie", framework.LoggingMiddleware(func(c *framework.Context) {
		cookie.ClearCookie(c.Response, "session")
		c.Response.Text(200, "Cookie deleted")
	}))

	server.Get("/", framework.LoggingMiddleware(func(c *framework.Context) {
		data := db.Select("users")

		c.JSON(200, data)
	}))

	server.Get("/{id}", func(c *framework.Context) {
		id := c.GetPathValue("id")

		data := db.SelectById("users", id)

		c.JSON(200, data)
	})

	server.Post("/", func(c *framework.Context) {
		data := Example{}

		c.BindBody(&data)

		db.Insert("users", data)

		c.JSON(200, data)
	})

	server.Put("/{id}", func(c *framework.Context) {
		id := c.GetPathValue("id")

		data := Example{}

		c.BindBody(&data)

		if ok := db.Update("users", id, data); !ok {
			c.Status(404)

			return
		}

		c.JSON(200, data)
	})

	server.Delete("/{id}", func(c *framework.Context) {
		id := c.GetPathValue("id")

		if ok := db.Delete("users", id); !ok {
			c.Status(404)

			return
		}

		c.Text(200, "Deleted")
	})

	server.Start(port)
}
