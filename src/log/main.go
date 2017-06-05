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

func LogFile() string {
    pathLog := beego.AppConfig.String("log::file")
    return pathLog
}

func main() {
	logFile := LogFile()
	logs.SetLogger(logs.AdapterFile,`{"filename":"`+logFile+`","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
	
	time.AfterFunc(10*time.Second,CallProducer)
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

func CallProducer(){
	_,err := httplib.Get("http://127.0.0.1:8080/v1/rabbitmq/producer").Debug(true).Response()
	if err != nil {log.Println(err)}
}