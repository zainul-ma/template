package routers

import (
	"customer/controllers"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSCond(func(ctx *context.Context) bool {
        // if ctx.Input.Domain() == "api.beego.me" {
        //     return true
        // }
        return true
    }),
		beego.NSBefore(Auth),
		beego.NSNamespace("/tbl_customer",
			beego.NSInclude(
				&controllers.TblCustomerController{},
			),
		),
		beego.NSAfter(AfterFunc),
	)
	beego.AddNamespace(ns)
}


// Auth to call authentication
func Auth(c *context.Context) {
	beego.Debug("checking.....")
}

// AfterFunc to execute progress after response
func AfterFunc(c *context.Context) {
	
}
