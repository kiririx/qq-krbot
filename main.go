package main

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"os"
	"qq-krbot/api"
	"qq-krbot/dao"
	"time"
)

func main() {
	mysqlUser := flag.String("mysql_user", "root", "")
	mysqlPass := flag.String("mysql_pass", "", "")
	mysqlHost := flag.String("mysql_host", "localhost", "")
	mysqlDb := flag.String("mysql_db", "qq_kr_bot", "")
	flag.Parse()
	dao.InitORM(fmt.Sprintf(
		`%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local&&timeout=1s&readTimeout=5s&writeTimeout=5s`,
		os.Getenv(*mysqlUser),
		os.Getenv(*mysqlPass),
		os.Getenv(*mysqlHost),
		os.Getenv(*mysqlDb),
	), 10, 500, time.Minute*15)
	r := gin.Default()
	r.Use(gin.Recovery())
	r.GET("/ping", api.Ping)
	r.POST("/api/bot", api.Bot)
	r.Use(static.Serve("/photo", static.LocalFile("./photo", true)))
	r.Run(":10013")
}

func setEnv() {

}
