package test

import (
  "testing"
  "customer/models"
  "log"
  "reflect"
  // "customer/struct"
  // "encoding/json"

  // "github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDataBase("default", "mysql", "root:@tcp(127.0.0.1:3306)/beego")
}

func TestGetAll(t *testing.T){
  var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

  l, _ := models.GetAllTblCustomer(query, fields, sortby, order, offset, limit)

  log.Println(reflect.TypeOf(l))
  log.Println(l);
  for key,val := range l {
    log.Println(key)
    log.Println(val)
  }
  t.Error("Expected")

}
