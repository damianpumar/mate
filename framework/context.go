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

func (c *Context) GetQueryParam(key string) string {
	return c.Request.GetQueryParam(key)
}

func (c *Context) BindBody(data interface{}) error {
	return c.Request.BindBody(data)
}

func (c *Context) Text(status int, text string) {
	c.Response.Text(status, text)
}

func (c *Context) JSON(status int, data interface{}) {
	c.Response.JSON(status, data)
}
