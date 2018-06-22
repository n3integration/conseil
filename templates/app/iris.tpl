package main

import (
    "github.com/kataras/iris"
)

var addr = iris.Addr("{{ .Host }}:{{ .Port }}")

func main() {
    // Create new router
    app := iris.New()

    // Register health endpoint
    app.Get("/health", health)

    // Now listening on: http://{{ .Host }}:{{ .Port }}
    // Application started. Press CTRL+C to shut down.
    app.Run(addr)
}

// Iris Handler
func health(ctx iris.Context) {
    ctx.JSON(iris.Map{
    	"status": "OK",
    })
}
