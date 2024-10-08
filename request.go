package mate

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

func (r *Request) GetPathValue(key string) string {
	return r.PathValue(key)
}

func (r *Request) BindBody(data interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return err
	}

	return nil
}
