package main

import (
	_ "log/routers"
	"time"
	"github.com/astaxie/beego/httplib"
	"log"

	"github.com/astaxie/beego"
)

func main() {
	time.AfterFunc(10*time.Second,CallFuncIndexTrendingPost2DaysAgo)
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

func CallFuncIndexTrendingPost2DaysAgo(){
	req,err := httplib.Get("http://127.0.0.1:8080/v1/rabbitmq/producer").Debug(true).Response()
	if err != nil {log.Println(err)}
	log.Println(req)
}
