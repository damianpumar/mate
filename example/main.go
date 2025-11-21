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

	// GET all users
	server.Get("/", mate.LoggingMiddleware(func(c *mate.Context) {
		data, err := db.Select("users")

		if err != nil {
			c.Error(500, err)
			return
		}

		c.JSON(200, data)
	}))

	// GET all users (typed version)
	server.Get("/users/typed", mate.LoggingMiddleware(func(c *mate.Context) {
		users, err := database.SelectTyped[Example](db, "users")

		if err != nil {
			c.Error(500, err)
			return
		}

		c.JSON(200, users)
	}))

	// GET user by ID
	server.Get("/{id}", func(c *mate.Context) {
		id := c.GetPathValue("id")

		data, err := db.SelectByID("users", id)

		if err != nil {
			if err == database.ErrRecordNotFound {
				c.Error(404, err)
				return
			}
			c.Error(500, err)
			return
		}

		c.JSON(200, data)
	})

	// GET user by ID (typed version)
	server.Get("/users/typed/{id}", func(c *mate.Context) {
		id := c.GetPathValue("id")

		user, err := database.SelectByIDTyped[Example](db, "users", id)

		if err != nil {
			if err == database.ErrRecordNotFound {
				c.Error(404, err)
				return
			}
			c.Error(500, err)
			return
		}

		c.JSON(200, user)
	})

	// POST create user
	server.Post("/", func(c *mate.Context) {
		data := Example{}

		if err := c.BindBody(&data); err != nil {
			c.Error(400, err)
			return
		}

		if err := db.Insert("users", data); err != nil {
			c.Error(500, err)
			return
		}

		c.JSON(201, data)
	})

	// PUT update user
	server.Put("/{id}", func(c *mate.Context) {
		id := c.GetPathValue("id")

		data := Example{}

		if err := c.BindBody(&data); err != nil {
			c.Error(400, err)
			return
		}

		if err := db.Update("users", id, data); err != nil {
			if err == database.ErrRecordNotFound {
				c.Status(404)
				return
			}
			c.Error(500, err)
			return
		}

		c.JSON(200, data)
	})

	// DELETE user
	server.Delete("/{id}", func(c *mate.Context) {
		id := c.GetPathValue("id")

		if err := db.Delete("users", id); err != nil {
			if err == database.ErrRecordNotFound {
				c.Status(404)
				return
			}
			c.Error(500, err)
			return
		}

		c.Text(200, "Deleted")
	})

	// GET users with filter
	server.Get("/users/filter", func(c *mate.Context) {
		name := c.Request.Request.URL.Query().Get("name")

		users, err := database.SelectWhereTyped(db, "users", func(u Example) bool {
			return u.Name == name
		})

		if err != nil {
			c.Error(500, err)
			return
		}

		c.JSON(200, users)
	})

	// GET count users
	server.Get("/users/count", func(c *mate.Context) {
		count, err := db.Count("users")

		if err != nil {
			c.Error(500, err)
			return
		}

		c.JSON(200, map[string]int{"count": count})
	})

	// GET check if user exists
	server.Get("/users/exists/{id}", func(c *mate.Context) {
		id := c.GetPathValue("id")

		exists, err := db.Exists("users", id)

		if err != nil {
			c.Error(500, err)
			return
		}

		c.JSON(200, map[string]bool{"exists": exists})
	})

	// DELETE all users (truncate)
	server.Delete("/users/truncate", func(c *mate.Context) {
		if err := db.Truncate("users"); err != nil {
			c.Error(500, err)
			return
		}

		c.Text(200, "All users deleted")
	})

	// POST bulk insert users
	server.Post("/users/bulk", func(c *mate.Context) {
		var users []Example

		if err := c.BindBody(&users); err != nil {
			c.Error(400, err)
			return
		}

		records := make([]interface{}, len(users))
		for i, u := range users {
			records[i] = u
		}

		if err := db.InsertMany("users", records); err != nil {
			c.Error(500, err)
			return
		}

		c.JSON(201, users)
	})

	// PUT update multiple users by name
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

		updated, err := database.UpdateWhereTyped(db, "users",
			func(u Example) bool {
				return u.Name == update.OldName
			},
			func(u Example) Example {
				u.Name = update.NewName
				return u
			})

		if err != nil {
			c.Error(500, err)
			return
		}

		c.JSON(200, map[string]interface{}{
			"updated": updated,
			"message": "Users updated successfully",
		})
	})

	// DELETE multiple users by name
	server.Delete("/users/by-name/{name}", func(c *mate.Context) {
		name := c.GetPathValue("name")

		deleted, err := database.DeleteWhereTyped(db, "users", func(u Example) bool {
			return u.Name == name
		})

		if err != nil {
			c.Error(500, err)
			return
		}

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
