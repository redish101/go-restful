package main

import (
	"fmt"
	"net/http"

	"github.com/redish101/go-restful"
)

func main() {
	container := restful.NewContainer()

	ws := restful.NewWebService().Path("/")

	ws.GET("/hello/{name}", helloWorld)

	container.Add(ws)

	http.ListenAndServe(":3000", container)
}

func helloWorld(req *restful.Request, resp *restful.Response) {
	// name := req.PathParameter("name")
	resp.WriteError(http.StatusInternalServerError, fmt.Errorf("aaa"))
}
