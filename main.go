package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"net/http"
	"qq-krbot/api"
)

func main() {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/api/bot", api.Bot)
	r.Use(static.Serve("/photo", static.LocalFile("./photo", true)))
	r.Run(":10011")
}
