package restful

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Container struct {
	mu       sync.Mutex
	mux      *http.ServeMux
	services []*WebService
	filters  []Filter
	produces []string
	consumes []string
	built    bool
}

func NewContainer() *Container {
	return &Container{
		mux: http.NewServeMux(),
	}
}

func (c *Container) Filter(f Filter) *Container {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.filters = append(c.filters, f)
	return c
}

func (c *Container) Produces(m ...string) *Container {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.produces = append(c.produces, m...)
	return c
}

func (c *Container) Consumes(m ...string) *Container {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.consumes = append(c.consumes, m...)
	return c
}

func (c *Container) Add(ws *WebService) *Container {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services = append(c.services, ws)
	c.built = false
	return c
}

func (c *Container) AddAll(services ...*WebService) *Container {
	for _, ws := range services {
		c.Add(ws)
	}
	return c
}

func (c *Container) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.build()

	ri := &responseInterceptor{ResponseWriter: w}
	c.mux.ServeHTTP(ri, r)

	if ri.intercepted {
		w.Header().Set("Content-Type", MIME_JSON+"; charset=utf-8")
		w.Header().Del("X-Content-Type-Options")

		msg := "not found"
		if ri.status == http.StatusMethodNotAllowed {
			msg = "method not allowed"
		}

		w.WriteHeader(ri.status)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": msg,
			"code":  ri.status,
		})
	}
}

func (c *Container) Handler() http.Handler {
	c.build()
	return c
}

func (c *Container) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, c.Handler())
}

func (c *Container) build() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.built {
		return
	}
	c.built = true

	patterns := make(map[string]bool)
	for _, ws := range c.services {
		for _, route := range ws.Routes() {
			pattern := route.pattern()
			if patterns[pattern] {
				log.Printf("[go-restful] WARNING: duplicate route pattern %q", pattern)
			}
			patterns[pattern] = true
			c.mux.HandleFunc(pattern, c.wrap(route))
		}
	}
}

func (c *Container) wrap(route *Route) http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := NewRequest(r)
		resp := NewResponse(w)
		route.Handler(req, resp)
	})

	filters := route.effectiveFilters(c)
	for i := len(filters) - 1; i >= 0; i-- {
		handler = filters[i](handler).(http.HandlerFunc)
	}

	produces := route.effectiveProduces(c)
	consumes := route.effectiveConsumes(c)

	return func(w http.ResponseWriter, r *http.Request) {
		if ri, ok := w.(*responseInterceptor); ok {
			ri.routeMatched = true
		}

		if len(consumes) > 0 && hasBody(r) {
			ct := r.Header.Get("Content-Type")
			if ct == "" || !acceptsMIME(ct, consumes) {
				w.Header().Set("Content-Type", MIME_JSON+"; charset=utf-8")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				_, _ = w.Write([]byte(`{"error":"unsupported media type","code":415}`))
				return
			}
		}

		if len(produces) > 0 {
			chosen := negotiate(r.Header.Get("Accept"), produces)
			if chosen == "" {
				w.Header().Set("Content-Type", MIME_JSON+"; charset=utf-8")
				w.WriteHeader(http.StatusNotAcceptable)
				_, _ = w.Write([]byte(`{"error":"not acceptable","code":406}`))
				return
			}
			w.Header().Set("Content-Type", chosen+"; charset=utf-8")
		}

		handler.ServeHTTP(w, r)
	}
}

func (c *Container) Services() []*WebService {
	return c.services
}

func hasBody(r *http.Request) bool {
	if r.ContentLength > 0 {
		return true
	}
	return len(r.TransferEncoding) > 0
}

type responseInterceptor struct {
	http.ResponseWriter
	routeMatched bool
	status       int
	statusSet    bool
	intercepted  bool
}

func (ri *responseInterceptor) WriteHeader(status int) {
	if ri.statusSet {
		return
	}
	ri.statusSet = true
	ri.status = status

	if !ri.routeMatched &&
		(status == http.StatusNotFound || status == http.StatusMethodNotAllowed) {
		ri.intercepted = true
		return
	}

	ri.ResponseWriter.WriteHeader(status)
}

func (ri *responseInterceptor) Write(b []byte) (int, error) {
	if ri.intercepted {
		return len(b), nil
	}
	return ri.ResponseWriter.Write(b)
}
