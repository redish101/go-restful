package restful

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
)

type Request struct {
	*http.Request
}

func NewRequest(r *http.Request) *Request {
	return &Request{Request: r}
}

func (r *Request) PathParameter(name string) string {
	return r.PathValue(name)
}

func (r *Request) QueryParameter(name string) string {
	return r.URL.Query().Get(name)
}

func (r *Request) HeaderParameter(name string) string {
	return r.Header.Get(name)
}

func (r *Request) ReadEntity(v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func (r *Request) ReadEntityWithContentType(v any) error {
	defer r.Body.Close()
	ct := r.Header.Get("Content-Type")
	if strings.HasPrefix(ct, MIME_XML) {
		return xml.NewDecoder(r.Body).Decode(v)
	}
	return json.NewDecoder(r.Body).Decode(v)
}

func (r *Request) ReadBody() ([]byte, error) {
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func (r *Request) SetContextValue(key, val any) *Request {
	r.Request = r.Request.WithContext(
		context.WithValue(r.Request.Context(), key, val),
	)
	return r
}

func (r *Request) ContextValue(key any) any {
	return r.Request.Context().Value(key)
}

func (r *Request) Unwrap() *http.Request {
	return r.Request
}
