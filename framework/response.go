package framework

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	http.ResponseWriter
}

func (r *Response) WithJson(status int, data interface{}) {
	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(status)
	json.NewEncoder(r).Encode(data)
}

func (r *Response) WithText(status int, message string) {
	r.Header().Set("Content-Type", "text/plain")
	r.WriteHeader(status)
	r.Write([]byte(message))
}
