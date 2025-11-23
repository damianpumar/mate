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

	server.Get("/users/typed", mate.LoggingMiddleware(func(c *mate.Context) {
		users := database.SelectTyped[Example](db, "users")
		c.JSON(200, users)
	}))

	server.Get("/{id}", func(c *mate.Context) {
		id := c.GetPathValue("id")
		data, found := db.SelectByID("users", id)

		if !found {
			c.Status(404)
			return
		}

		c.JSON(200, data)
	})

	server.Get("/users/typed/{id}", func(c *mate.Context) {
		id := c.GetPathValue("id")
		user, found := database.SelectByIDTyped[Example](db, "users", id)

		if !found {
			c.Status(404)
			return
		}

		c.JSON(200, user)
	})

	server.Post("/", func(c *mate.Context) {
		data := Example{}

		if err := c.BindBody(&data); err != nil {
			c.Error(400, err)
			return
		}

		db.Insert("users", data)
		c.JSON(201, data)
	})

	server.Put("/{id}", func(c *mate.Context) {
		id := c.GetPathValue("id")
		data := Example{}

		if err := c.BindBody(&data); err != nil {
			c.Error(400, err)
			return
		}

		if !db.Update("users", id, data) {
			c.Status(404)
			return
		}

		c.JSON(200, data)
	})

	server.Delete("/{id}", func(c *mate.Context) {
		id := c.GetPathValue("id")

		if !db.Delete("users", id) {
			c.Status(404)
			return
		}

		c.Text(200, "Deleted")
	})

	server.Get("/users/filter", func(c *mate.Context) {
		name := c.Request.Request.URL.Query().Get("name")
		users := database.SelectWhereTyped(db, "users", func(u Example) bool {
			return u.Name == name
		})

		c.JSON(200, users)
	})

	server.Get("/users/count", func(c *mate.Context) {
		count := db.Count("users")
		c.JSON(200, map[string]int{"count": count})
	})

	server.Get("/users/exists/{id}", func(c *mate.Context) {
		id := c.GetPathValue("id")
		exists := db.Exists("users", id)
		c.JSON(200, map[string]bool{"exists": exists})
	})

	server.Delete("/users/truncate", func(c *mate.Context) {
		db.Truncate("users")
		c.Text(200, "All users deleted")
	})

	server.Post("/users/bulk", func(c *mate.Context) {
		var users []Example

		if err := c.BindBody(&users); err != nil {
			c.Error(400, err)
			return
		}

		database.InsertManyTyped(db, "users", users)
		c.JSON(201, users)
	})

	server.Put("/users/bulk-update", func(c *mate.Context) {
		type BulkUpdate struct {
			OldName string `json:"old_name"`
			NewName string `json:"new_name"`
		}

		var update BulkUpdate
		if err := c.BindBody(&update); err != nil {
			c.Error(400, err)
			return
		}

		updated := database.UpdateWhereTyped(db, "users",
			func(u Example) bool {
				return u.Name == update.OldName
			},
			func(u Example) Example {
				u.Name = update.NewName
				return u
			})

		c.JSON(200, map[string]interface{}{
			"updated": updated,
			"message": "Users updated successfully",
		})
	})

	server.Delete("/users/by-name/{name}", func(c *mate.Context) {
		name := c.GetPathValue("name")
		deleted := database.DeleteWhereTyped(db, "users", func(u Example) bool {
			return u.Name == name
		})

		c.JSON(200, map[string]interface{}{
			"deleted": deleted,
			"message": "Users deleted successfully",
		})
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
