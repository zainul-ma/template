package test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/astaxie/beego"
	"github.com/streadway/amqp"

	"customer/thirdparty"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(
		filepath.Join(file, "../../"+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

func TestRabbitCredentialAndConnection(t *testing.T) {
	mqURL := thirdparty.CredMq()

	_, err := amqp.Dial(mqURL)

	_, errMq := thirdparty.ConnectMq(mqURL)

	Convey("Subject: RabbitMq Test\n", t, func() {
		Convey("Rabbit mq credential", func() {
			So(mqURL, ShouldEqual, beego.AppConfig.String("mq::local"))
		})
		Convey("Rabbit connecction Url", func() {
			So(err, ShouldEqual, nil)
		})
		Convey("Rabbit Mq connection third party", func() {
			So(errMq, ShouldEqual, nil)
		})
	})
}
