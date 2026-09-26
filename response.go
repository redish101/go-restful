package restful

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func NewResponse(w http.ResponseWriter) *Response {
	return &Response{ResponseWriter: w}
}

func (r *Response) WriteEntity(v any) error {
	return r.WriteHeaderAndEntity(http.StatusOK, v)
}

func (r *Response) WriteHeaderAndEntity(status int, v any) error {
	r.Header().Set("Content-Type", MIME_JSON+"; charset=utf-8")
	r.WriteHeader(status)
	return json.NewEncoder(r).Encode(v)
}

func (r *Response) WriteEntityWithContentType(v any, accept string) error {
	if strings.Contains(accept, MIME_XML) {
		r.Header().Set("Content-Type", MIME_XML+"; charset=utf-8")
		r.WriteHeader(http.StatusOK)
		_, err := io.WriteString(r, xml.Header)
		if err != nil {
			return err
		}
		return xml.NewEncoder(r).Encode(v)
	}
	return r.WriteEntity(v)
}

func (r *Response) WriteError(status int, err error) error {
	return r.WriteHeaderAndEntity(status, map[string]any{
		"error": err.Error(),
	})
}

func (r *Response) WriteErrorString(status int, msg string) error {
	return r.WriteHeaderAndEntity(status, map[string]any{
		"error": msg,
	})
}

func (r *Response) WriteText(status int, text string) error {
	r.Header().Set("Content-Type", MIME_TEXT+"; charset=utf-8")
	r.WriteHeader(status)
	_, err := io.WriteString(r, text)
	return err
}

func (r *Response) WriteHeader(status int) {
	if r.written {
		return
	}
	r.statusCode = status
	r.written = true
	r.ResponseWriter.WriteHeader(status)
}

func (r *Response) StatusCode() int {
	return r.statusCode
}

func (r *Response) Written() bool {
	return r.written
}

func (r *Response) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}
