package test

import (
	"customer/models"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/astaxie/beego"

	mgo "gopkg.in/mgo.v2"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(
		filepath.Join(file, "../../"+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

func TestDBCredential(t *testing.T) {
	dbURL := models.DbCred()

	fmt.Println(dbURL)

	Convey("Subject: Monggo Credential\n", t, func() {
		Convey("Monggo db config url", func() {
			So(dbURL, ShouldEqual, beego.AppConfig.String("mongodb::local"))
		})
	})
}

func TestConnection(t *testing.T) {
	dbURL := models.DbCred()

	_, err := mgo.Dial(dbURL)

	session := models.ConnectMongo()

	defer session.Close()

	Convey("Subject: Monggo Connection\n", t, func() {
		Convey("Monggo conecction error should be nil", func() {
			So(err, ShouldEqual, nil)
		})
		Convey("Monggo session mongotonic mode", func() {
			So(session.Mode(), ShouldEqual, mgo.Monotonic)
		})
	})
}
