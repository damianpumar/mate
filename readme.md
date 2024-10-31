# Mate 🧉

Mate is a simple, lightweight, fast an zero-dependency web server with batteries included http server.

# Getting Started 🚀

## Installation 📦

```bash
go get github.com/damianpumar/mate
```

## Usage 🏃‍♀️

```go
package main

import (
	"github.com/damianpumar/mate"
	"github.com/damianpumar/mate/database"
)

type User struct {
	Name string `json:"name"`
}

func main() {
	server := mate.New()
	db := database.Connect() // Optional json database

	server.Get("/", func(ctx *mate.Context) {
		ctx.JSON(200, db.Select("users"))
	})

	server.Get("/users", func(ctx *mate.Context) {
		db.Insert("users", &User{Name: "Damian"})

		ctx.Text(200, "Inserted")
	})

  	server.Post("/", func(c *mate.Context) {
		data := User{}

		c.BindBody(&data)

		db.Insert("users", data)

		c.JSON(200, data)
	})

	server.Start("3000")
}
```

## Running the server 🚀

```bash
go run main.go
```

## License 📄

MIT
