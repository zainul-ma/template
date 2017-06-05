package models

import(
	"github.com/astaxie/beego"

	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"

	"os"
)

var envOs = os.Getenv("GOENV")
type TypeDbCred struct {
	Local string `json:"local"`
	Dev string `json:"dev"`
	Prod string `json:"prod"`
}

func DbCred() string {
	db_url := ""
    if envOs == "local" {
    	db_url = beego.AppConfig.String("mongodb:local")
	}else if envOs == "dev" {
		db_url = beego.AppConfig.String("mongodb:dev")
	}else if envOs == "prod" {
		db_url = beego.AppConfig.String("mongodb:prod")
	}

    return db_url
}

func ConnectMongo() *mgo.Session {
	dbUrl := DbCred();
	session,err := mgo.Dial(dbUrl)
	CheckErr(err,"error connect mongo DB")

	if envOs != "prod" {
		mgo.SetDebug(true)
	}
	session.SetMode(mgo.Monotonic,true)

	return session
}

func CheckErr(err error, msg string) {
	if err != nil {
		beego.Warning(msg)
   		beego.Critical(err)
	}
}

