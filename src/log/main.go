package main

import (
	_ "log/routers"
	"time"
	"github.com/astaxie/beego/httplib"
	"log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	// "os"
)

func main() {
	// aa := os.Getenv("GOENV")

	logs.SetLogger(logs.AdapterFile,`{"filename":"log/project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
	time.AfterFunc(10*time.Second,CallProducer)
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

func CallProducer(){
	req,err := httplib.Get("http://127.0.0.1:8080/v1/rabbitmq/producer").Debug(true).Response()
	if err != nil {log.Println(err)}
	log.Println(req)
}
