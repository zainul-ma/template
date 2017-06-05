package helper

import (
	"github.com/astaxie/beego/context"
	"encoding/json"
	"github.com/astaxie/beego"
)

// HeaderAll to get all header request
func HeaderAll(c *context.Context) string {
	headerAll := make(map[string]string)
	for k,v := range c.Request.Header {
		headerAll[k] = v[0]
	}
	strJSON,err := json.Marshal(headerAll)
	if err != nil {
		beego.Warning("error Marshal header request")
		beego.Error(err)
	}
	return string(strJSON)
}