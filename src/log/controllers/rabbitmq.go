package controllers

import (
	"log/models"
	// "encoding/json"

	"github.com/astaxie/beego"
)

// Operations about object
type RabbitmqController struct {
	beego.Controller
}

func (o *RabbitmqController) GetMQ() {
  models.Receiver()
  o.ServeJSON();
}
