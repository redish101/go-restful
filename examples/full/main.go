package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redish101/go-restful"
)

type Article struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateArticleRequest struct {
	Title string `json:"title" validate:"required,min=1,max=200"`
	Body  string `json:"body"  validate:"required,max=10000"`
}

var (
	mu       sync.RWMutex
	articles = map[string]*Article{}
)

func listArticles(req *restful.Request, resp *restful.Response) {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]*Article, 0, len(articles))
	for _, a := range articles {
		out = append(out, a)
	}
	_ = resp.WriteEntity(map[string]any{"items": out, "count": len(out)})
}

func getArticle(req *restful.Request, resp *restful.Response) {
	mu.RLock()
	a, ok := articles[req.PathParameter("id")]
	mu.RUnlock()
	if !ok {
		_ = resp.WriteErrorString(http.StatusNotFound, "article not found")
		return
	}
	_ = resp.WriteEntity(a)
}

func createArticle(req *restful.Request, resp *restful.Response) {
	var body CreateArticleRequest
	if err := req.BindAndValidate(&body); err != nil {
		_ = resp.WriteError(http.StatusUnprocessableEntity, err)
		return
	}
	a := &Article{
		ID:        uuid.NewString(),
		Title:     body.Title,
		Body:      body.Body,
		CreatedAt: time.Now(),
	}
	mu.Lock()
	articles[a.ID] = a
	mu.Unlock()

	resp.Header().Set("Location", "/articles/"+a.ID)
	_ = resp.WriteHeaderAndEntity(http.StatusCreated, a)
}

func deleteArticle(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")
	mu.Lock()
	_, ok := articles[id]
	delete(articles, id)
	mu.Unlock()
	if !ok {
		_ = resp.WriteErrorString(http.StatusNotFound, "article not found")
		return
	}
	resp.WriteHeader(http.StatusNoContent)
}

func main() {
	restful.SetDefaultProduces(restful.MIME_JSON)
	restful.SetDefaultConsumes(restful.MIME_JSON)

	c := restful.NewContainer()
	c.Produces(restful.MIME_JSON)
	c.Consumes(restful.MIME_JSON)
	c.Filter(restful.Recover())
	c.Filter(restful.RequestID())
	c.Filter(restful.Logger())
	c.Filter(restful.CORS())

	ws := restful.NewWebService().Path("/articles")
	ws.Filter(restful.RequireJSON())

	ws.GET("", listArticles).WithDoc("List all articles")
	ws.GET("/{id}", getArticle).WithDoc("Get an article by ID")
	ws.POST("", createArticle).
		WithDoc("Create a new article").
		WithConsumes(restful.MIME_JSON)
	ws.DELETE("/{id}", deleteArticle).WithDoc("Delete an article")

	c.Add(ws)

	health := restful.NewWebService().Path("/health")
	health.GET("", func(req *restful.Request, resp *restful.Response) {
		_ = resp.WriteEntity(map[string]any{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})
	c.Add(health)

	http.ListenAndServe(":3000", c)
}
