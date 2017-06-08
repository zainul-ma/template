package test

import (
	"bytes"
	"customer/models"
	structCustomer "customer/structs"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	_ "customer/routers"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
)

func removeAll() {
	connection := models.ConnectMongo()
	defer connection.Close()
	c := connection.DB(models.GetDBName()).C(models.GetDBName())
	c.RemoveAll(nil)
}

func buildRecord() (bson.ObjectId, structCustomer.Customer) {
	ID := bson.NewObjectId()

	v := structCustomer.Customer{
		ID:       ID,
		Fullname: "fake",
		Email:    "fake@faker.com",
		Username: "fakerusername",
	}

	return ID, v
}

func init() {
	removeAll()
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(
		filepath.Join(file, "../../"+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestCustomerGet is a sample to run an endpoint test
func TestCustomerGet(t *testing.T) {
	removeAll()

	for i := 0; i < 11; i++ {
		_, v := buildRecord()

		if err := models.AddCustomer(&v); err != nil {
			t.Error("-")
		}
	}

	r, _ := http.NewRequest("GET",
		"/v1/tbl_customer"+
			"?sortby=email,fullname"+
			"&order=asc,desc"+
			"&fields=fullname,email"+
			"&limit=10",
		nil)

	w := httptest.NewRecorder()

	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestGetAllCustomer", "Code[%d]\n%s", w.Code,
		w.Body.String())

	Convey("Subject: Customer Endpoint Get All\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}

//TestCustomerGetById
func TestCustomerGetById(t *testing.T) {
	removeAll()

	ID, v := buildRecord()

	if err := models.AddCustomer(&v); err != nil {
		t.Error("-")
	}

	r, _ := http.NewRequest("GET",
		"/v1/tbl_customer/"+ID.Hex(),
		nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestCustomerGetById", "Code[%d]\n%s", w.Code,
		w.Body.String())

	Convey("Subject: Customer Endpoint Get By ID\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}

//
func TestCustomerDelete(t *testing.T) {
	removeAll()

	ID, v := buildRecord()

	if err := models.AddCustomer(&v); err != nil {
		t.Error("-")
	}

	r, _ := http.NewRequest("DELETE", "/v1/tbl_customer/"+ID.Hex(), nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestDeleteCustomer", "Code[%d]\n%s", w.Code,
		w.Body.String())

	Convey("Subject: Customer Endpoint Delete\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
			So(w.Body.String(), ShouldEqual, "\"OK\"")
		})
	})
}

//
func TestCustomerUpdate(t *testing.T) {
	removeAll()

	ID, v := buildRecord()

	if err := models.AddCustomer(&v); err != nil {
		t.Error("-")
	}

	byteJson, _ := json.Marshal(v)

	r, _ := http.NewRequest("PUT", "/v1/tbl_customer/"+ID.Hex(),
		bytes.NewBuffer(byteJson))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestUpdateCustomer", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Customer Endpoint Update\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
	})
}

func TestCustomerAdd(t *testing.T) {
	removeAll()

	_, v := buildRecord()

	byteJson, _ := json.Marshal(v)

	r, _ := http.NewRequest("POST", "/v1/tbl_customer",
		bytes.NewBuffer(byteJson))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestGet", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Customer Endpoint Add\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 201)
		})

		byteResult, _ := json.Marshal(v)

		body := strings.Replace(w.Body.String(), "\n", "", -1)

		body = strings.Replace(body, " ", "", -1)

		Convey("Response should be have response for created", func() {
			So(body, ShouldResemble, string(byteResult))
		})

		Convey("Decode json like Struct type", func() {
			value := &structCustomer.Customer{}
			err := json.Unmarshal([]byte(body), value)
			So(err, ShouldEqual, nil)
			So(value.Email, ShouldEqual, "fake@faker.com")
		})
	})
}
