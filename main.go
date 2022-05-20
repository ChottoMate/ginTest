package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", func(context *gin.Context) {
		context.HTML(200, "index.html", gin.H{})
	})

	router.Run()
}
