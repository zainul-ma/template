package main

import (
	_ "customer/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func LogFile() string {
    pathLog := beego.AppConfig.String("log::file")
    return pathLog
}

func init() {
	logFile := LogFile()
	logs.SetLogger(logs.AdapterFile,`{"filename":"`+logFile+`","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
