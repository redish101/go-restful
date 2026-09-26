package restful

import "net/http"

type HandlerFunc func(req *Request, resp *Response)

type Filter func(http.Handler) http.Handler

type Handler = http.Handler

type EntityReader interface {
	ReadEntity(v any) error
}

type EntityWriter interface {
	WriteEntity(v any) error
	WriteHeaderAndEntity(status int, v any) error
}

type Validatable interface {
	Validate() error
}
