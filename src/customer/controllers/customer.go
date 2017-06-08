package controllers

import (
	"encoding/json"
	"errors"
	"math/rand"
	"strconv"
	"strings"

	"customer/helper"
	models "customer/models"
	structCustomer "customer/structs"
	"customer/thirdparty"

	"github.com/astaxie/beego"
	// "github.com/astaxie/beego/context"
)

// TblCustomerController operations for TblCustomer
type TblCustomerController struct {
	beego.Controller
}

// TrackingOutputCustomer method
func TrackingOutputCustomer(c *TblCustomerController) {
	SendMq(c)
	c.ServeJSON()
}

// SendMq method
func SendMq(c *TblCustomerController) {
	fromService := beego.BConfig.AppName
	reqBody := c.Ctx.Input.CopyBody(int64(1200))
	resBody, _ := json.Marshal(c.Data["json"])
	headerAll := helper.HeaderAll(c.Ctx)
	toService := ""

	reqID := c.Ctx.Input.Header("reqID")
	newRequest := false
	if reqID == "" {
		newRequest = true
		PtrReqID(&reqID, rand.Int(), &fromService, "client", &toService,
			beego.BConfig.AppName)
	} else {

	}

	thirdparty.SendMq(reqBody, fromService, toService, headerAll, reqID,
		newRequest, "req")
	thirdparty.SendMq(resBody, fromService, toService, headerAll, reqID,
		newRequest, "res")
}

// PtrReqID method
func PtrReqID(reqID *string, val int, fromService *string,
	valFromService string, toService *string, valToService string) {
	*reqID = strconv.Itoa(val)
	*fromService = valFromService
	*toService = valToService
}

// URLMapping ...
func (c *TblCustomerController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create TblCustomer
// @Param	body		body 	models.TblCustomer	true		"body for TblCustomer content"
// @Success 201 {int} models.TblCustomer
// @Failure 403 body is empty
// @router / [post]
func (c *TblCustomerController) Post() {
	var v structCustomer.Customer
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err2 := models.AddCustomer(&v); err2 == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	TrackingOutputCustomer(c)
}

// GetOne ...
// @Title Get One
// @Description get TblCustomer by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.TblCustomer
// @Failure 403 :id is empty
// @router /:id [get]
func (c *TblCustomerController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")

	v, err := models.GetCustomerByID(idStr)

	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}

	TrackingOutputCustomer(c)
}

// GetAll ...
// @Title Get All
// @Description get TblCustomer
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	latestID	query	string	false	"Start position of result set. Must be an ID"
// @Success 200 {object} models.TblCustomer
// @Failure 403
// @router / [get]
func (c *TblCustomerController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int = 10
	var latestID string

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt("limit"); err == nil {
		limit = v
	}
	// latestID: 0 (default is 0)
	if v := c.GetString("latestID"); v != "" {
		latestID = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllCustomer(query, fields, sortby, order, latestID,
		limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}

	TrackingOutputCustomer(c)
}

// Put ...
// @Title Put
// @Description update the TblCustomer
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.TblCustomer	true		"body for TblCustomer content"
// @Success 200 {object} models.TblCustomer
// @Failure 403 :id is not int
// @router /:id [put]
func (c *TblCustomerController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	var v structCustomer.Customer
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err2 := models.UpdateCustomer(idStr, &v); err2 == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}

	TrackingOutputCustomer(c)
}

// Delete ...
// @Title Delete
// @Description delete the TblCustomer
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *TblCustomerController) Delete() {
	idStr := c.Ctx.Input.Param(":id")

	if err := models.DeleteCustomer(idStr); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}

	TrackingOutputCustomer(c)
}
