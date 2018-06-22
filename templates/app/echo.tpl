package main

import (
    "net/http"

    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
)

var addr = "{{ .Host }}:{{ .Port }}"

func main() {
    // Create new router
    r := echo.New()

    // Setup common middleware
    r.Use(
        middleware.Logger(),
        middleware.Recover(),
    )

    // Register health endpoint
    r.GET("/health", health)

    // Now listening on: http://{{ .Host }}:{{ .Port }}
    // Application started. Press CTRL+C to shut down.
    r.Logger.Fatal(r.Start(addr))
}

// Echo handler
func health(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]string{
        "status": "OK",
    })
}