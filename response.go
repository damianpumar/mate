package mate

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type Response struct {
	http.ResponseWriter
}

func (r *Response) JSON(status int, data interface{}) {
	r.Header().Set("Content-Type", "application/json")

	r.WriteHeader(status)

	json.NewEncoder(r).Encode(data)
}

func (r *Response) Text(status int, message string) {
	r.Header().Set("Content-Type", "text/plain")

	r.WriteHeader(status)

	r.Write([]byte(message))
}

func (r *Response) Status(status int) {
	r.WriteHeader(status)
}

func (r *Response) Error(status int, err error) {
	r.Header().Set("Content-Type", "text/plain")

	r.WriteHeader(status)

	r.Write([]byte(err.Error()))
}

func (r *Response) Render(status int, file string, data interface{}) {
	r.Header().Set("Content-Type", "text/html")
	r.WriteHeader(status)

	html := template.Must(template.ParseFiles(file))

	if err := html.Execute(r, data); err != nil {
		http.Error(r, err.Error(), http.StatusInternalServerError)
	}
}
