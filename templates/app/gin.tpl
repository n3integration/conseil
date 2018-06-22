package main

import (
    "github.com/gin-gonic/gin"
)

var addr = "{{ .Host }}:{{ .Port }}"

func main() {
    // Create new router
    r := gin.Default()

    // Register health endpoint
    r.GET("/health", health)

    // Now listening on: http://{{ .Host }}:{{ .Port }}
    // Application started. Press CTRL+C to shut down.
    r.Run(addr)
}

// Gin handler
func health(c *gin.Context) {
    c.JSON(200, gin.H{
        "status": "OK",
    })
}