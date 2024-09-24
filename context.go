package mate

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

func (c *Context) GetPathValue(key string) string {
	return c.Request.GetPathValue(key)
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

func (c *Context) Status(status int) {
	c.Response.Status(status)
}

func (c *Context) Error(status int, err error) {
	c.Response.Error(status, err)
}

func (c *Context) Render(status int, template string, data interface{}) {
	c.Response.Render(status, template, data)
}
