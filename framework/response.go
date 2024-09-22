package framework

import (
	"encoding/json"
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
