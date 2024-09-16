package framework

import "net/http"

type Request struct {
	*http.Request
}

func (r *Request) GetQueryParam(key string) string {
	values := r.URL.Query()

	return values.Get(key)
}
