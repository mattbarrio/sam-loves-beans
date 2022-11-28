package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// for some reason gin doesn't let you pass params to functions inside request handlers, so I had to hack this middleware
func OpenAIConfig() gin.HandlerFunc {
	godotenv.Load()

	var apiKey = os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	return func(c *gin.Context) {
		c.Set("apiKey", apiKey)
		c.Next()
	}

}

func main() {

	// start statsWorker to track ip requests
	go statsWorker()

	r := gin.New()
	r.Use(OpenAIConfig())
	// use stats and crash recovery middlewares
	r.Use(rateLimit, gin.Recovery(), gin.Logger())
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static/")
	r.GET("/", index)
	r.GET("/ping", ping)
	r.GET("/api/beans/text", newCompletionRequest)
	r.GET("/api/beans/image", newImageRequest)

	// m := autocert.Manager{
	// 	Prompt:     autocert.AcceptTOS,
	// 	HostPolicy: autocert.HostWhitelist("sam-loves-beans.com"),
	// 	Cache:      autocert.DirCache("/var/www/.cache"),
	// }
	r.Run(":8080")
	//log.Fatal(autotls.RunWithManager(r, &m))

}
