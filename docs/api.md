# go-restful api

## API 文档

### Container

```go
type Container struct { /* ... */ }

func NewContainer() *Container
```

| 方法 | 说明 |
|---|---|
| `Filter(f Filter) *Container` | 追加容器级过滤器（最外层） |
| `Produces(m ...string) *Container` | 设置容器级默认响应 MIME |
| `Consumes(m ...string) *Container` | 设置容器级默认请求 MIME |
| `Add(ws *WebService) *Container` | 注册一个 WebService |
| `AddAll(services ...*WebService) *Container` | 批量注册 |
| `ServeHTTP(w, r)` | 实现 `http.Handler` |
| `Handler() http.Handler` | 返回构建后的 `http.Handler` |
| `ListenAndServe(addr string) error` | 一行启动 |
| `Services() []*WebService` | 返回所有 WebService |

**示例**

```go
c := restful.NewContainer()
c.Produces(restful.MIME_JSON)
c.Consumes(restful.MIME_JSON)
c.Filter(restful.Recover())
c.Filter(restful.Logger())
c.Add(usersWS).Add(articlesWS)

// 塞进 http.Server
srv := &http.Server{Addr: ":8080", Handler: c.Handler()}

// 或挂到别的 router 下
mux := http.NewServeMux()
mux.Handle("/api/", http.StripPrefix("/api", c.Handler()))
```

---

### WebService

```go
func NewWebService() *WebService
```

| 方法 | 说明 |
|---|---|
| `Path(p string) *WebService` | 设置根路径，如 `/users` |
| `Filter(f Filter) *WebService` | 追加服务级过滤器 |
| `Produces(m ...string) *WebService` | 服务级默认响应 MIME |
| `Consumes(m ...string) *WebService` | 服务级默认请求 MIME |
| `GET/POST/PUT/PATCH/DELETE/HEAD/OPTIONS(path, handler) *Route` | 注册路由 |
| `Route(method, path string, handler) *Route` | 自定义方法 |
| `Routes() []*Route` | 返回所有路由 |
| `RootPath() string` | 返回根路径 |

**示例**

```go
ws := restful.NewWebService().Path("/articles")
ws.Produces(restful.MIME_JSON)
ws.Filter(restful.RequireJSON())

ws.GET("", listArticles)              // GET    /articles
ws.GET("/{id}", getArticle)           // GET    /articles/{id}
ws.POST("", createArticle)            // POST   /articles
ws.PUT("/{id}", updateArticle)        // PUT    /articles/{id}
ws.DELETE("/{id}", deleteArticle)     // DELETE /articles/{id}

c.Add(ws)
```

---

### Route

由 `ws.GET(...)` 等返回，可链式配置：

| 方法 | 说明 |
|---|---|
| `WithFilter(f Filter) *Route` | 追加路由级过滤器 |
| `WithProduces(m ...string) *Route` | 声明响应 MIME |
| `WithConsumes(m ...string) *Route` | 声明请求 MIME |
| `WithDoc(doc string) *Route` | 附加文档 |
| `WithParameter(p *Parameter) *Route` | 附加参数元数据 |

字段：`Method`、`Path`、`Handler`、`Filters`、`Produces`、`Consumes`、`Doc`、`Parameters`。

**示例**

```go
ws.POST("", createArticle).
    WithConsumes(restful.MIME_JSON).
    WithProduces(restful.MIME_JSON).
    WithFilter(authMiddleware).
    WithDoc("Create a new article").
    WithParameter(&restful.Parameter{
        Name:        "X-Api-Key",
        Description: "API key",
        DataType:    "string",
        Required:    true,
    })
```

---

### Request

`Request` 内嵌 `*http.Request`，**所有标准库字段和方法都可用**。

```go
type Request struct {
    *http.Request
}
```

#### 便捷方法

| 方法 | 说明 |
|---|---|
| `PathParameter(name string) string` | 路径参数，基于 `r.PathValue` |
| `QueryParameter(name string) string` | 查询参数 |
| `HeaderParameter(name string) string` | 请求头 |
| `ReadEntity(v any) error` | 解码 JSON 请求体 |
| `ReadEntityWithContentType(v any) error` | 按 `Content-Type` 选择解码器（JSON/XML） |
| `ReadBody() ([]byte, error)` | 读取原始请求体 |
| `BindAndValidate(v any) error` | **解码 + 校验**（推荐） |
| `SetContextValue(key, val any) *Request` | 写入 request context |
| `ContextValue(key any) any` | 读取 request context |
| `Unwrap() *http.Request` | 返回底层 `*http.Request` |

#### 标准库透传

```go
req.Method              // "GET"
req.URL.Path            // "/articles/42"
req.Header.Get("X-Foo") // 请求头
req.Context()           // context.Context
req.WithContext(ctx)    // 派生新请求
req.RemoteAddr          // 客户端地址
```

**示例**

```go
func handler(req *restful.Request, resp *restful.Response) {
    id := req.PathParameter("id")
    q  := req.QueryParameter("q")
    ua := req.Header.Get("User-Agent")

    req = req.SetContextValue("trace_id", "abc123")
    _ = req.ContextValue("trace_id")

    someStdFunc(req.Unwrap())
}
```

---

### Response

`Response` 内嵌 `http.ResponseWriter`，**所有标准库字段和方法都可用**。

```go
type Response struct {
    http.ResponseWriter
}
```

#### 便捷方法

| 方法 | 说明 |
|---|---|
| `WriteEntity(v any) error` | 写 JSON，状态码 200 |
| `WriteHeaderAndEntity(status int, v any) error` | 写 JSON，指定状态码 |
| `WriteEntityWithContentType(v any, accept string) error` | 按 `Accept` 选择 JSON/XML |
| `WriteError(status int, err error) error` | 统一错误格式 `{"error":..,"code":..}` |
| `WriteErrorString(status int, msg string) error` | 同上，字符串版 |
| `WriteText(status int, text string) error` | 写纯文本 |
| `StatusCode() int` | 已写入的状态码，未写入返回 0 |
| `Written() bool` | 是否已写入响应 |
| `Unwrap() http.ResponseWriter` | 返回底层 `http.ResponseWriter` |

#### 标准库透传

```go
resp.Header().Set("X-Total-Count", "42")
resp.WriteHeader(http.StatusNoContent)
resp.Write([]byte("raw"))
```

**示例**

```go
func handler(req *restful.Request, resp *restful.Response) {
    _ = resp.WriteEntity(user)
    _ = resp.WriteHeaderAndEntity(http.StatusCreated, user)
    _ = resp.WriteError(http.StatusNotFound, errors.New("not found"))

    resp.Header().Set("Location", "/users/42")
    resp.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(resp).Encode(user)
}
```

---

### Filter

```go
type Filter func(http.Handler) http.Handler
```

**就是标准库中间件签名**，与 `alice` / `chi/middleware` / `gorilla/handlers` 完全兼容。

执行顺序：**Container → WebService → Route**（外层先执行）。

**示例**

```go
func Timing(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s took %s", r.URL.Path, time.Since(start))
    })
}

c.Filter(Timing)                     // 全局
ws.Filter(Timing)                    // 服务级
ws.GET("/x", h).WithFilter(Timing)   // 路由级
```

---

### 校验

内置 [go-playground/validator/v10](https://github.com/go-playground/validator) 集成。

#### 基础用法

```go
type CreateUserRequest struct {
    Name  string `json:"name"  validate:"required,min=2,max=50"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age"   validate:"gte=0,lte=150"`
}

func createUser(req *restful.Request, resp *restful.Response) {
    var body CreateUserRequest
    if err := req.BindAndValidate(&body); err != nil {
        _ = resp.WriteError(http.StatusUnprocessableEntity, err)
        return
    }
    // ...
}
```

失败响应：

```json
{
  "error": "field 'Name' failed on 'required'; field 'Email' failed on 'required'",
  "code": 422
}
```

#### 自定义校验

实现 `Validatable` 接口，**优先级高于 tag 校验**：

```go
func (r *UpdateUserRequest) Validate() error {
    if r.Name == nil && r.Email == nil {
        return errors.New("nothing to update")
    }
    return nil
}
```

#### 错误类型

```go
type ValidateError struct {
    Field   string
    Tag     string
    Value   any
    Message string
}

type ValidateErrors []*ValidateError
```

可按需 `errors.As` 拆解：

```go
if err := req.BindAndValidate(&body); err != nil {
    var verrs restful.ValidateErrors
    if errors.As(err, &verrs) {
        for _, e := range verrs {
            log.Printf("%s: %s", e.Field, e.Message)
        }
    }
    _ = resp.WriteError(http.StatusUnprocessableEntity, err)
}
```

#### 全局校验器

```go
v := restful.GetValidator()
v.RegisterValidation("mytag", func(fl validator.FieldLevel) bool { ... })

restful.SetValidator(customValidator)

err := restful.ValidateStruct(&user)
```

#### 常用 tag 速查

| Tag | 说明 |
|---|---|
| `required` | 必填 |
| `email` | 邮箱格式 |
| `url` | URL 格式 |
| `min` / `max` | 长度或数值范围 |
| `len` | 精确长度 |
| `gte` / `lte` | 数值范围 |
| `oneof` | 枚举（`oneof=red green blue`） |
| `omitempty` | 空值跳过（指针、切片等） |

完整列表见 [validator 文档](https://pkg.go.dev/github.com/go-playground/validator/v10)。

---

## 默认 MIME 与回退顺序

`Produces` / `Consumes` 支持**四级回退**，就近覆盖：

```
Route  >  WebService  >  Container  >  包级默认
```

| 层级 | 声明方式 | 优先级 | 作用范围 |
|---|---|---|---|
| **Route** | `ws.GET(...).WithProduces(...)` | 最高 | 单条路由 |
| **WebService** | `ws.Produces(...)` | 次高 | 该服务的所有路由 |
| **Container** | `c.Produces(...)` | 中 | 该容器下所有 WebService |
| **包级默认** | `restful.SetDefaultProduces(...)` | 最低 | 进程内所有路由 |
| **无** | — | — | 不设置 Content-Type，由 handler 自行决定 |

**包级默认**

```go
restful.SetDefaultProduces(restful.MIME_JSON)
restful.SetDefaultConsumes(restful.MIME_JSON)

// 读取
restful.DefaultProduces()
restful.DefaultConsumes()
```

---

## 内容协商

### 请求侧

若声明了 `Consumes` 且请求携带 body：

- 校验 `Content-Type` 是否在 `Consumes` 列表中；
- 不匹配返回 `415 Unsupported Media Type`；
- 兼容 `chunked` 传输（`TransferEncoding` 非空即视为有 body）。

### 响应侧

若声明了 `Produces`：

- 解析 `Accept` 头（含 `q=` 权重、`*/*`、`type/*` 通配）；
- 从候选中挑选最佳匹配，写入 `Content-Type`；
- 无法匹配返回 `406 Not Acceptable`；
- `Accept` 为空时取 `Produces` 的第一项。

```bash
# 匹配成功
curl -H 'Accept: application/json' localhost:8080/users
# → 200, Content-Type: application/json

# 通配符
curl -H 'Accept: */*' localhost:8080/users

# 不匹配
curl -H 'Accept: image/png' localhost:8080/users
# → 406 {"error":"not acceptable","code":406}

# 非法请求体 MIME
curl -X POST localhost:8080/users -H 'Content-Type: text/plain' -d 'x'
# → 415 {"error":"unsupported media type","code":415}
```

---

## 内置过滤器

| 过滤器 | 说明 |
|---|---|
| `Recover()` | 捕获 panic，返回 500 JSON |
| `Logger()` | 访问日志（方法、路径、耗时） |
| `RequestID()` | 注入 `X-Request-ID` 响应头，写入 context |
| `CORS(origins ...string)` | 跨域支持，默认 `*`，自动处理 `OPTIONS` |
| `RequireJSON()` | POST/PUT/PATCH 强制 `application/json` |

**示例**

```go
c.Filter(restful.Recover())
c.Filter(restful.RequestID())
c.Filter(restful.Logger())
c.Filter(restful.CORS("https://example.com"))
```

**上下文工具函数**

```go
req := restful.SetUser(r, "alice")
user := restful.GetUser(r)
rid := restful.GetRequestID(r)
```

---

## 路由语法

基于 Go 1.22+ `http.ServeMux`：

| 模式 | 说明 |
|---|---|
| `/users` | 精确匹配 |
| `/users/{id}` | 单段路径参数，`req.PathParameter("id")` 获取 |
| `/files/{path...}` | 匹配剩余所有路径 |
| `/users/{id}/posts/{postID}` | 多参数 |
| `{$}` | 仅匹配根路径 `/` |

**注意**：`ServeMux` 对冲突路由会 panic，例如 `/users/{id}` 与 `/users/{name}` 不能同时注册。

---

## 中间件生态

因为 `Filter` 就是标准库签名，可以直接使用第三方中间件。

**alice**

```go
import "github.com/justinas/alice"

chain := alice.New(restful.Recover, restful.Logger)
c.Filter(func(next http.Handler) http.Handler {
    return chain.Then(next)
})
```

**chi/middleware**

```go
import chimw "github.com/go-chi/chi/v5/middleware"

c.Filter(chimw.RequestID)
c.Filter(chimw.RealIP)
c.Filter(chimw.Compress(5))
c.Filter(chimw.Timeout(10 * time.Second))
```

**gorilla/handlers**

```go
import "github.com/gorilla/handlers"

c.Filter(handlers.CompressHandler)
c.Filter(handlers.RecoveryHandler())
```

---

## 错误处理

`Response` 提供两个统一错误输出方法：

```go
// 结构化
_ = resp.WriteError(http.StatusNotFound, errors.New("user not found"))
// → {"error":"user not found","code":404}

// 字符串
_ = resp.WriteErrorString(http.StatusBadRequest, "invalid input")
// → {"error":"invalid input","code":400}
```

推荐在 handler 里集中处理：

```go
func (h *Handler) Get(req *restful.Request, resp *restful.Response) {
    user, err := h.store.Get(req.PathParameter("id"))
    if err != nil {
        switch {
        case errors.Is(err, ErrNotFound):
            _ = resp.WriteErrorString(http.StatusNotFound, "not found")
        default:
            _ = resp.WriteError(http.StatusInternalServerError, err)
        }
        return
    }
    _ = resp.WriteEntity(user)
}
```

---

