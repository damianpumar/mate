package main

import (
	"minimal-http-server/database"
	"minimal-http-server/films"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		films := films.GetFilms()

		return c.JSON(films)
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
