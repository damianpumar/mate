package framework

import (
	"net/http"
)

type Context struct {
	Response *Response
	Request  *Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Response: &Response{w},
		Request:  &Request{r},
	}
}
