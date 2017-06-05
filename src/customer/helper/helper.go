package helper

import (
	"github.com/astaxie/beego/context"
	"encoding/json"
)

func HeaderAll(c *context.Context) string {
	headerAll := make(map[string]string)
	for k,v := range c.Request.Header {
		headerAll[k] = v[0]
	}
	strJson,_ := json.Marshal(headerAll)
	return string(strJson)
}