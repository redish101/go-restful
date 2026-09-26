package restful

import (
	"net/http"
	"strings"
)

type Route struct {
	Method     string
	Path       string
	Handler    HandlerFunc
	Filters    []Filter
	Produces   []string
	Consumes   []string
	Doc        string
	Parameters []*Parameter

	ws *WebService
}

type Parameter struct {
	Name        string
	Description string
	DataType    string
	Required    bool
	Default     string
}

func (r *Route) WithFilter(f Filter) *Route {
	r.Filters = append(r.Filters, f)
	return r
}

func (r *Route) WithProduces(m ...string) *Route {
	r.Produces = append(r.Produces, m...)
	return r
}

func (r *Route) WithConsumes(m ...string) *Route {
	r.Consumes = append(r.Consumes, m...)
	return r
}

func (r *Route) WithDoc(doc string) *Route {
	r.Doc = doc
	return r
}

func (r *Route) WithParameter(p *Parameter) *Route {
	r.Parameters = append(r.Parameters, p)
	return r
}

func (r *Route) fullPath() string {
	root := ""
	if r.ws != nil {
		root = r.ws.rootPath
	}
	return joinPath(root, r.Path)
}

func (r *Route) pattern() string {
	return r.Method + " " + r.fullPath()
}

func (r *Route) effectiveFilters(container *Container) []Filter {
	var filters []Filter
	filters = append(filters, container.filters...)
	if r.ws != nil {
		filters = append(filters, r.ws.filters...)
	}
	filters = append(filters, r.Filters...)
	return filters
}

func (r *Route) effectiveConsumes(c *Container) []string {
	if len(r.Consumes) > 0 {
		return r.Consumes
	}
	if r.ws != nil && len(r.ws.consumes) > 0 {
		return r.ws.consumes
	}
	if c != nil && len(c.consumes) > 0 {
		return c.consumes
	}
	return defaultConsumes
}

func (r *Route) effectiveProduces(c *Container) []string {
	if len(r.Produces) > 0 {
		return r.Produces
	}
	if r.ws != nil && len(r.ws.produces) > 0 {
		return r.ws.produces
	}
	if c != nil && len(c.produces) > 0 {
		return c.produces
	}
	return defaultProduces
}

func joinPath(root, sub string) string {
	root = strings.TrimRight(root, "/")
	sub = strings.TrimLeft(sub, "/")

	if root == "" && sub == "" {
		return "/"
	}
	if root == "" {
		return "/" + sub
	}
	if sub == "" {
		return root
	}
	return root + "/" + sub
}

func acceptsMIME(header string, allowed []string) bool {
	if header == "" {
		return false
	}
	ct := strings.TrimSpace(strings.Split(header, ";")[0])
	for _, a := range allowed {
		if strings.EqualFold(ct, a) {
			return true
		}
	}
	return false
}

func httpMethod(method string) string {
	return strings.ToUpper(strings.TrimSpace(method))
}

var _ = http.MethodGet
