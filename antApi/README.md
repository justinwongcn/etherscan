# Ant Web 框架使用指南

## 快速开始

### 安装

```bash
go get github.com/justinwongcn/ant
```

### 基本用法

以下是一个简单的示例，展示了如何创建一个 HTTP 服务器并处理请求：

```go
package main

import (
    "fmt"
    "github.com/justinwongcn/ant"
)

func main() {
    // 创建一个新的 HTTP 服务器
    server := ant.NewHTTPServer()

    // 注册路由处理函数
    server.Handle("GET /hello", func(ctx *ant.Context) {
        // 获取查询参数
        nameVal := ctx.QueryValue("name")
        name, err := nameVal.String()
        if err != nil || name == "" {
            name = "World"
        }

        // 设置响应头
        ctx.Resp.Header().Set("Content-Type", "text/plain; charset=utf-8")

        // 返回响应
        fmt.Fprintf(ctx.Resp, "Hello, %s!", name)
    })

    // 启动服务器
    if err := server.Run(":8080"); err != nil {
        fmt.Printf("Server failed to start: %v\n", err)
    }
}
```

## 核心功能

### 路由注册

框架支持简单的路由注册方式：

```go
// 基本路由
server.Handle("GET /path", handler)

// 带参数的路由（需要 Go 1.22+）
server.Handle("GET /users/{id}", func(ctx *ant.Context) {
    id := ctx.Req.PathValue("id")
    // 处理请求...
})
```

### 请求参数处理

```go
// 获取查询参数
nameVal := ctx.QueryValue("name")
name, err := nameVal.String()

// 获取路径参数（需要 Go 1.22+）
id := ctx.Req.PathValue("id")
```

### 中间件支持

框架支持中间件机制，可以在请求处理前后添加自定义逻辑：

```go
// 定义中间件
logger := func(next ant.HandleFunc) ant.HandleFunc {
    return func(ctx *ant.Context) {
        // 请求前的处理
        fmt.Println("Before request")
        
        next(ctx)
        
        // 请求后的处理
        fmt.Println("After request")
    }
}

// 注册中间件
server.Use(logger)
```

### 文件处理

框架提供了文件上传和静态资源服务的支持：

```go
// 文件上传处理
uploader := &ant.FileUploader{
    FileField: "file",
    DstPathFunc: func(fh *multipart.FileHeader) string {
        return filepath.Join("uploads", fh.Filename)
    },
}
server.Handle("POST /upload", uploader.Handle())

// 静态资源服务
static := ant.NewStaticResourceHandler("public", "/static/")
server.Handle("GET /static/{file}", static.Handle)
```

### 模板渲染

框架支持 Go 模板引擎：

```go
engine := &ant.GoTemplateEngine{}

// 从文件加载模板
engine.LoadFromFiles("templates/*.html")

// 渲染模板
result, err := engine.Render(context.Background(), "template.html", data)
```

## 最佳实践

1. 使用中间件处理通用逻辑，如日志记录、错误处理等
2. 合理组织路由，使用有意义的 URL 路径
3. 适当处理错误，返回合适的状态码和错误信息
4. 使用适当的响应格式，设置正确的 Content-Type

## 示例

访问示例应用：

```
http://localhost:8080/hello           # 返回 "Hello, World!"
http://localhost:8080/hello?name=John # 返回 "Hello, John!"
```