package main

import (
	"newapi-mini/middleware"
	"newapi-mini/model"
	"newapi-mini/relay"

	"github.com/gin-gonic/gin"
)

func main() {
	model.InitDB()

	r := gin.Default()
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "ok",
		})
	})

	r.POST("/v1/chat/completions", middleware.AuthMiddleware(), func(ctx *gin.Context) {
		relay.RelayChat(ctx)
	})

	r.Run(":8080")
}
