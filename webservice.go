package restful

import (
	"strings"
)

type WebService struct {
	rootPath string
	routes   []*Route
	filters  []Filter
	consumes []string
	produces []string
}

func NewWebService() *WebService {
	return &WebService{}
}

func (ws *WebService) Path(p string) *WebService {
	ws.rootPath = strings.TrimRight(p, "/")
	if ws.rootPath == "" {
		ws.rootPath = "/"
	}
	return ws
}

func (ws *WebService) Filter(f Filter) *WebService {
	ws.filters = append(ws.filters, f)
	return ws
}

func (ws *WebService) Consumes(m ...string) *WebService {
	ws.consumes = append(ws.consumes, m...)
	return ws
}

func (ws *WebService) Produces(m ...string) *WebService {
	ws.produces = append(ws.produces, m...)
	return ws
}

func (ws *WebService) GET(path string, handler HandlerFunc) *Route {
	return ws.add(httpMethod("GET"), path, handler)
}

func (ws *WebService) POST(path string, handler HandlerFunc) *Route {
	return ws.add(httpMethod("POST"), path, handler)
}

func (ws *WebService) PUT(path string, handler HandlerFunc) *Route {
	return ws.add(httpMethod("PUT"), path, handler)
}

func (ws *WebService) PATCH(path string, handler HandlerFunc) *Route {
	return ws.add(httpMethod("PATCH"), path, handler)
}

func (ws *WebService) DELETE(path string, handler HandlerFunc) *Route {
	return ws.add(httpMethod("DELETE"), path, handler)
}

func (ws *WebService) HEAD(path string, handler HandlerFunc) *Route {
	return ws.add(httpMethod("HEAD"), path, handler)
}

func (ws *WebService) OPTIONS(path string, handler HandlerFunc) *Route {
	return ws.add(httpMethod("OPTIONS"), path, handler)
}

func (ws *WebService) Route(method, path string, handler HandlerFunc) *Route {
	return ws.add(httpMethod(method), path, handler)
}

func (ws *WebService) add(method, path string, handler HandlerFunc) *Route {
	r := &Route{
		Method:  method,
		Path:    path,
		Handler: handler,
		ws:      ws,
	}
	ws.routes = append(ws.routes, r)
	return r
}

func (ws *WebService) Routes() []*Route {
	return ws.routes
}

func (ws *WebService) RootPath() string {
	return ws.rootPath
}
