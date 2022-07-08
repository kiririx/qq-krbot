package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"qq-krbot/api"
	"qq-krbot/env"
	"qq-krbot/work"
)

func main() {
	new(work.CollectWorker).Start()
	new(work.SendWorker).Start()
	r := gin.Default()
	r.Use(gin.Recovery())
	r.GET("/ping", api.Ping)
	r.POST("/api/bot", api.Bot)
	// r.Use(static.Serve("/photo", static.LocalFile("./photo", true)))
	r.GET("/photo/:tag/:path", func(c *gin.Context) {
		uri, err := url.QueryUnescape(c.Request.RequestURI)
		if err != nil {
			return
		}
		b, err := ioutil.ReadFile("." + uri)
		if err != nil {
			return
		}
		c.Data(http.StatusOK, "image/jpeg", b)
	})
	r.Run(":" + env.Conf["serve.port"])

}
