package main

import (
	_ "customer/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/logs"
)

type TypeGlobalCred struct {
	Local string `json:"local"`
	Dev string `json:"dev"`
	Prod string `json:"prod"`
}

func(this *TypeGlobalCred) LogFile() {

}

func(this *TypeGlobalCred) DbCred() {
	
}

func init() {
	logs.SetLogger(logs.AdapterFile,`{"filename":"log/project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
	orm.RegisterDataBase("default", "mysql", "root:@tcp(127.0.0.1:3306)/beego")
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
