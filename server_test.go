package mate_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/damianpumar/mate"
)

type ServerTest struct {
	Serve *mate.Server
}

func createTestFramework() *ServerTest {
	server := mate.New()

	return &ServerTest{Serve: server}
}

func (s *ServerTest) get(path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	s.Serve.ServeHTTP(res, req)

	return res
}

func (s *ServerTest) post(path string, body interface{}) *httptest.ResponseRecorder {
	bodyParsed, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", path, bytes.NewBuffer(bodyParsed))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	s.Serve.ServeHTTP(res, req)

	return res
}

func (s *ServerTest) assertEqual(t *testing.T, res *httptest.ResponseRecorder, expected string) {
	if res.Code != 200 {
		t.Errorf("Expected status code 200, got %d", res.Code)
	}

	var expectedJSON, actualJSON interface{}

	err := json.Unmarshal([]byte(expected), &expectedJSON)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %s", err)
	}

	err = json.Unmarshal(res.Body.Bytes(), &actualJSON)
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

func TestServer(t *testing.T) {
	type Example struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	testServer := createTestFramework()

	testServer.Serve.Get("/", func(c *mate.Context) {
		c.JSON(200, []Example{{Id: "1", Name: "John Doe"}})
	})

	testServer.Serve.Get("/{id}", func(c *mate.Context) {
		id := c.GetPathValue("id")

		c.JSON(200, Example{Id: id, Name: "John Doe"})
	})

	testServer.Serve.Post("/", func(c *mate.Context) {
		data := Example{}

		c.BindBody(&data)

		c.JSON(200, data)
	})

	t.Run("GET /", func(t *testing.T) {
		res := testServer.get("/")

		testServer.assertEqual(t, res, `[{"id":"1","name":"John Doe"}]`)
	})

	t.Run("POST /", func(t *testing.T) {
		res := testServer.post("/", Example{Id: "1", Name: "John Doe"})

		testServer.assertEqual(t, res, `{"id":"1","name":"John Doe"}`)
	})

	t.Run("GET /{id}", func(t *testing.T) {
		res := testServer.get("/20")

		testServer.assertEqual(t, res, `{"id":"20","name":"John Doe"}`)
	})
}
