package framework

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	*http.Request
}

func (r *Request) GetQueryParam(key string) string {
	values := r.URL.Query()

	return values.Get(key)
}

func (r *Request) ParseBody(data interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return err
	}

	return nil
}
