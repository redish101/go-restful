# go-restful

> 基于标准库 `net/http.ServeMux` 的更现代、简洁的对 [emicklei/go-restful](https://github.com/emicklei/go-restful) 的实现。

## 安装

```bash
go get github.com/redish101/go-restful
```


## 快速开始

```go
package main

import (
    "log"

    restful "github.com/redish101/go-restful"
)

type User struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

func main() {
    c := restful.NewContainer()
    c.Filter(restful.Recover())
    c.Filter(restful.Logger())
    c.Produces(restful.MIME_JSON)

    ws := restful.NewWebService().Path("/users")
    ws.GET("", listUsers)
    ws.GET("/{id}", getUser)

    c.Add(ws)
    log.Fatal(c.ListenAndServe(":8080"))
}

func listUsers(req *restful.Request, resp *restful.Response) {
    _ = resp.WriteEntity([]User{{ID: "1", Name: "Alice"}})
}

func getUser(req *restful.Request, resp *restful.Response) {
    _ = resp.WriteEntity(User{ID: req.PathParameter("id"), Name: "Alice"})
}
```

## 许可

MIT License © redish101
