package main

import (
	"flag"
	"time"

	"github.com/damianpumar/mate"
	"github.com/damianpumar/mate/database"
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

	server := mate.New()

	db := database.Connect()

	cookie := mate.NewSecureCookie("my-secret")

	server.Get("/cookie", mate.LoggingMiddleware(func(c *mate.Context) {
		if err := cookie.SetEncryptedCookie(c.Response, "session", "user12345", 30*time.Second); err != nil {

			c.Error(500, err)

			return
		}

		c.Text(200, "Cookie set")
	}))

	server.Get("/read-cookie", mate.LoggingMiddleware(func(c *mate.Context) {
		value, err := cookie.GetEncryptedCookie(c.Request.Request, "session")

		if err != nil {
			c.Error(500, err)

			return
		}

		c.Text(200, value)
	}))

	server.Get("/delete-cookie", mate.LoggingMiddleware(func(c *mate.Context) {
		cookie.ClearCookie(c.Response, "session")
		c.Response.Text(200, "Cookie deleted")
	}))

	server.Get("/", mate.LoggingMiddleware(func(c *mate.Context) {
		data := db.Select("users")

		c.JSON(200, data)
	}))

	server.Get("/{id}", func(c *mate.Context) {
		id := c.GetPathValue("id")

		data := db.SelectById("users", id)

		c.JSON(200, data)
	})

	server.Post("/", func(c *mate.Context) {
		data := Example{}

		c.BindBody(&data)

		db.Insert("users", data)

		c.JSON(200, data)
	})

	server.Put("/{id}", func(c *mate.Context) {
		id := c.GetPathValue("id")

		data := Example{}

		c.BindBody(&data)

		if ok := db.Update("users", id, data); !ok {
			c.Status(404)

			return
		}

		c.JSON(200, data)
	})

	server.Delete("/{id}", func(c *mate.Context) {
		id := c.GetPathValue("id")

		if ok := db.Delete("users", id); !ok {
			c.Status(404)

			return
		}

		c.Text(200, "Deleted")
	})

	server.Group("/hello", func(g *mate.Group) {
		g.Get("/world", func(c *mate.Context) {
			c.Text(200, "Hello, World!")
		})
	})

	server.Get("/template", func(c *mate.Context) {
		c.Render(200, "index.html", &Example{
			Id:   "1",
			Name: "Damián",
		})
	})

	server.Folder("/static", "./static")
	server.File("/json", "./static/hello.json")

	server.Start(*port)
}
