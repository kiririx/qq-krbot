package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"qq-krbot/api"
	"qq-krbot/work"
)

func main() {
	new(work.Collect).Start()
	r := gin.Default()
	r.Use(gin.Recovery())
	r.GET("/ping", api.Ping)
	r.POST("/api/bot", api.Bot)
	r.Use(static.Serve("/photo", static.LocalFile("./photo", true)))
	r.Run(":10013")

}
