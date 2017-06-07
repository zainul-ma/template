package test

import (
	"customer/models"
	structCustomer "customer/structs"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func removeAll() {
	connection := models.ConnectMongo()

	defer connection.Close()

	c := connection.DB("").C("customer")

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
}

func TestAddCustomer(t *testing.T) {
	removeAll()

	ID, v := buildRecord()

	if err := models.AddCustomer(&v); err != nil {
		t.Error("Should have created customer value")
	}

	v, err := models.GetCustomerByID(ID.Hex())

	if err != nil {
		t.Error("Should have created customer by ID")
	}

	if v.ID != ID {
		t.Error("Should have same ID , from initial ID and created ID")
	}
}

func TestGetAllCustomer(t *testing.T) {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int = 10
	var offset string

	removeAll()

	for i := 0; i < 11; i++ {
		_, v := buildRecord()

		if err := models.AddCustomer(&v); err != nil {
			t.Error("-")
		}
	}

	l, _ := models.GetAllCustomer(query, fields, sortby, order, offset, limit)

	if len(l) == 0 {
		t.Error("Should have records")
	}

	if len(l) != 10 {
		t.Error("Should be have 10 records")
	}
}

func TestGetCustomerByID(t *testing.T) {
	removeAll()

	ID, v := buildRecord()

	if err := models.AddCustomer(&v); err != nil {
		t.Error("Should have created customer value")
	}

	value, err := models.GetCustomerByID(ID.Hex())

	if err != nil {
		t.Error("Shoudl have get customer by ID")
	}

	if v != value {
		t.Error("Should be same value initial and get by ID")
	}
}

func TestGetDeletedByID(t *testing.T) {
	removeAll()

	ID, v := buildRecord()

	if err := models.AddCustomer(&v); err != nil {
		t.Error("Should have created customer value")
	}

	if err := models.DeleteCustomer(ID.Hex()); err != nil {
		t.Error("Should not have customer with ID created")
	}

	value, err := models.GetCustomerByID(ID.Hex())

	if err == nil {
		t.Error("Should not have error when get record by ID")
	}

	if value == v {
		t.Error("Should not have a record have beend deleted")
	}
}

func TestUpdateCustomer(t *testing.T) {
	removeAll()

	ID, v := buildRecord()

	u := structCustomer.Customer{
		Fullname: "ronyUpdated",
		Email:    "rony@fakerUpdated.com",
		Username: "ronyrusernameUpdated",
	}

	if err := models.AddCustomer(&v); err != nil {
		t.Error("Should have created customer value")
	}

	if err := models.UpdateCustomer(ID.Hex(), &u); err != nil {
		t.Error("Should have updated customer")
	}

	value, err := models.GetCustomerByID(ID.Hex())

	if err != nil {
		t.Error("Should have created or updated customer by ID")
	}

	if value.ID != ID {
		t.Error("Should have same ID , from initial ID and created ID")
	}

	u.ID = ID

	if value != u {
		t.Error("Should have same value, from updated and the select updated")
	}
}
