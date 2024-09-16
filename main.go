package main

import (
	"encoding/json"
	"flag"
	"minimal-http-server/database"
	"minimal-http-server/films"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

type User struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

var (
	version = flag.Bool("version", false, "prints the version")
	metrics = flag.Bool("metrics", false, "enable /metrics endpoint")
)

func main() {
	flag.Parse()

	if *version {
		println("v1.0.0")
		return
	}

	app := fiber.New()

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: encryptcookie.GenerateKey(),
	}))

	if *metrics {
		app.Get("/metrics", monitor.New(monitor.Config{Title: "Web server metrics"}))
	}

	app.Get("/", func(c *fiber.Ctx) error {
		films := films.GetFilms()

		return c.JSON(films)
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		sample := &User{
			Username: "admin",
			Role:     "admin",
		}

		sampleJSON, _ := json.Marshal(sample)

		c.Cookie(&fiber.Cookie{
			Name:     "session",
			Value:    string(sampleJSON),
			HTTPOnly: true,
			Expires:  time.Now().Add(10 * time.Second),
		})

		return c.JSON(fiber.Map{
			"message": "Logged in",
		})
	})

	app.Get("/logout", func(c *fiber.Ctx) error {
		c.ClearCookie("session")

		return c.JSON(fiber.Map{
			"message": "Logged out",
		})
	})

	app.Get("/cookie", func(c *fiber.Ctx) error {
		user := c.Cookies("session")

		if user == "" {
			return c.Status(403).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		return c.JSON(user)
	})

	app.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		film, err := films.GetFilm(id)
		if err != nil {

			return c.Status(404).JSON(fiber.Map{
				"error": "Film not found",
			})
		}

		return c.JSON(film)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		film := database.Film{
			Title:    c.FormValue("title"),
			Director: c.FormValue("director"),
		}

		films.AddFilm(film)

		return c.JSON(film)
	})

	app.Listen(":3000")
}
