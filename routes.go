package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func rateLimit(c *gin.Context) {
	if strings.Contains(c.FullPath(), "api") || c.FullPath() == "" {
		ip := c.ClientIP()
		value := int(ips.Add(ip, 1))
		if value >= 3 {
			fmt.Printf("rate limit: ip blocked - %s, count - %d\n", ip, value)
			c.Abort()
			c.String(http.StatusTooManyRequests, "Try again later")
		}
	}
}

// The index function loads the index template, redners responses if there are any, and add a js event listener to the button to grab more bean content async
func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":     "Beans!",
		"responses": getResponses(),
	})
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
