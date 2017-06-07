package models

import(
	"github.com/astaxie/beego"

	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"

	"os"
)

var envOs = os.Getenv("GOENV")
// DbCred call DB Cred
func DbCred() string {
	dbURL := ""
    if envOs == "local" {
    	dbURL = beego.AppConfig.String("mongodb::local")
	}else if envOs == "dev" {
		dbURL = beego.AppConfig.String("mongodb::dev")
	}else if envOs == "prod" {
		dbURL = beego.AppConfig.String("mongodb::prod")
	}

    return dbURL
}

// ConnectMongo to connect Mongo DB
func ConnectMongo() *mgo.Session {
	dbURL := DbCred();
	session,err := mgo.Dial(dbURL)
	CheckErr(err,"error connect mongo DB")

	if envOs != "prod" {
		mgo.SetDebug(true)
	}
	session.SetMode(mgo.Monotonic,true)

	return session
}

// CheckErr to check Error
func CheckErr(err error, msg string) {
	if err != nil {
		beego.Warning(msg)
   		beego.Critical(err)
	}
}

