package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"minimal/database"
	"minimal/framework"
	fwk "minimal/http"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Example struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func TestServer(t *testing.T) {

	os.Mkdir("database", 0755)
	defer os.RemoveAll("database")

	server := fwk.New()
	db := database.Connect()

	server.Get("/", func(c *framework.Context) {
		data := db.Select("users")

		c.JSON(200, data)
	})

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

	t.Run("POST /", func(t *testing.T) {
		body, _ := json.Marshal(Example{Id: "1", Name: "John Doe"})
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.JSONEq(t, `{"id":"1","name":"John Doe"}`, w.Body.String())
	})

	t.Run("GET /", func(t *testing.T) {
		body, _ := json.Marshal(Example{Id: "1", Name: "John Doe"})
		post, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
		post.Header.Set("Content-Type", "application/json")
		postW := httptest.NewRecorder()

		server.ServeHTTP(postW, post)

		get, _ := http.NewRequest("GET", "/", nil)
		get.Header.Set("Content-Type", "application/json")

		getW := httptest.NewRecorder()
		server.ServeHTTP(getW, get)

		assert.Equal(t, 200, getW.Code)
		fmt.Println(getW.Body.String())
		assert.JSONEq(t, `[{"id":"1","name":"John Doe"}]`, getW.Body.String())
	})

	t.Run("GET /{id}", func(t *testing.T) {
		body, _ := json.Marshal(Example{Id: "1", Name: "John Doe"})
		post, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
		post.Header.Set("Content-Type", "application/json")
		postW := httptest.NewRecorder()

		server.ServeHTTP(postW, post)

		get, _ := http.NewRequest("GET", "/1", nil)
		get.Header.Set("Content-Type", "application/json")

		getW := httptest.NewRecorder()
		server.ServeHTTP(getW, get)

		assert.Equal(t, 200, getW.Code)
		fmt.Println(getW.Body.String())
		assert.JSONEq(t, `{"id":"1","name":"John Doe"}`, getW.Body.String())
	})
}
