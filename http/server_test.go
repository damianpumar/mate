package http_test

import (
	"bytes"
	"encoding/json"
	"minimal/database"
	"minimal/framework"
	h "minimal/http"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

type Example struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func TestServer(t *testing.T) {

	server := h.New()
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
		os.Mkdir("database", 0755)
		defer os.RemoveAll("database")

		body, _ := json.Marshal(Example{Id: "1", Name: "John Doe"})
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		assertEqual(t, w, `{"id":"1","name":"John Doe"}`)
	})

	t.Run("GET /", func(t *testing.T) {
		os.Mkdir("database", 0755)
		defer os.RemoveAll("database")

		body, _ := json.Marshal(Example{Id: "1", Name: "John Doe"})
		post, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
		post.Header.Set("Content-Type", "application/json")
		postW := httptest.NewRecorder()

		server.ServeHTTP(postW, post)

		get, _ := http.NewRequest("GET", "/", nil)
		get.Header.Set("Content-Type", "application/json")

		getW := httptest.NewRecorder()
		server.ServeHTTP(getW, get)

		assertEqual(t, getW, `[{"id":"1","name":"John Doe"}]`)
	})

	t.Run("GET /{id}", func(t *testing.T) {
		os.Mkdir("database", 0755)
		defer os.RemoveAll("database")

		body, _ := json.Marshal(Example{Id: "1", Name: "John Doe"})
		post, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
		post.Header.Set("Content-Type", "application/json")
		postW := httptest.NewRecorder()

		server.ServeHTTP(postW, post)

		get, _ := http.NewRequest("GET", "/1", nil)
		get.Header.Set("Content-Type", "application/json")

		getW := httptest.NewRecorder()
		server.ServeHTTP(getW, get)

		assertEqual(t, getW, `{"id":"1","name":"John Doe"}`)
	})
}

func assertEqual(t *testing.T, rec *httptest.ResponseRecorder, expected string) {
	if rec.Code != 200 {
		t.Errorf("Expected status code 200, got %d", rec.Code)
	}

	var expectedJSON, actualJSON interface{}

	err := json.Unmarshal([]byte(expected), &expectedJSON)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %s", err)
	}

	err = json.Unmarshal(rec.Body.Bytes(), &actualJSON)
	if err != nil {
		t.Fatalf("Failed to unmarshal actual JSON: %s", err)
	}

	// Print the actual and expected JSON for debugging
	t.Logf("Expected JSON: %+v", expectedJSON)
	t.Logf("Actual JSON: %+v", actualJSON)

	if !reflect.DeepEqual(expectedJSON, actualJSON) {
		t.Errorf("Expected JSON %+v, got %+v", expectedJSON, actualJSON)
	}
}
