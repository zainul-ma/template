// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"customer/controllers"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego"
	// "customer/thirdparty"
	"log"
	// "math/rand"
	// "strconv"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSCond(func(ctx *context.Context) bool {
        if ctx.Input.Domain() == "api.beego.me" {
            return true
        }
        return true
    }),
		beego.NSBefore(Auth),
		beego.NSNamespace("/tbl_customer",
			beego.NSInclude(
				&controllers.TblCustomerController{},
			),
		),
	)
	beego.AddNamespace(ns)
}


func Auth(c *context.Context){
	log.Println("")
	// SendMQ(c)
	/*
	Input Request Data [Body Header]
	x := int64(120)
	t := c.Input.CopyBody(x)
	log.Println(string(t))

	l := c.Input.Header("key")
	log.Println(l)
	*/




	beego.Debug("checking.....")
	// c.Output.Body([]byte("bob"))
}

// func SendMQ(c *context.Context){
// 	reqID := ""
// 	fromService := beego.BConfig.AppName
// 	inputReqBody := c.Input.CopyBody(int64(1200))
// 	headerAll := c.Input.HeaderAll()
//
// 	if c.Input.Header("reqID") == "" {
// 		RandomID := strconv.Itoa(rand.Int())
// 		reqID = RandomID
// 		thirdparty.SendMQ(inputReqBody,fromService,"",headerAll,reqID)
// 	}
// }
