# go-restful

> 基于标准库 `net/http.ServeMux` 的更现代、简洁的对 [emicklei/go-restful](https://github.com/emicklei/go-restful) 的实现。

## 安装

```bash
go get github.com/redish101/go-restful
```

---

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

---

## 核心概念

```
Container                     ← 顶层容器，持有 http.ServeMux、全局过滤器与默认 MIME
  └── WebService[]            ← 一组共享根路径的路由（如 /users）
        └── Route[]           ← 单个「方法 + 路径 + handler」
              └── Filter[]    ← 路由级过滤器
```

- **Container**：一个 `http.Handler`，负责把路由注册到 `ServeMux`，串接过滤器链，提供默认 MIME。
- **WebService**：按资源划分路由，共享 `Path` 前缀与默认 MIME。
- **Route**：单条路由，可挂路由级过滤器、声明 `Produces`/`Consumes`、附加文档元数据。
- **Filter**：标准库风格中间件，三层可叠加，执行顺序为 `Container → WebService → Route`。

---

## 许可

MIT License © redish101
