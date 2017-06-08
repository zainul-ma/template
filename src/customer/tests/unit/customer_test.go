package test

import (
	"customer/models"
	structCustomer "customer/structs"
	"testing"

	"gopkg.in/mgo.v2/bson"

	. "github.com/smartystreets/goconvey/convey"
)

func removeAll() {
	connection := models.ConnectMongo()
	defer connection.Close()
	c := connection.DB(models.GetDBName()).C(models.GetTableName())
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

	errAdd := models.AddCustomer(&v)

	vGetByID, err := models.GetCustomerByID(ID.Hex())

	Convey("Subject: Unit Test Adding Customer\n", t, func() {
		Convey("Should have success created customer", func() {
			So(errAdd, ShouldEqual, nil)
		})
		Convey("Should have success get customer by ID", func() {
			So(err, ShouldEqual, nil)
		})
		Convey("Should have get created customer by ID", func() {
			So(v, ShouldResemble, vGetByID)
		})
		Convey("Should have same ID , from initial ID and created ID", func() {
			So(vGetByID.ID, ShouldEqual, ID)
		})
	})
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

	Convey("Subject: Unit Test Get Customer\n", t, func() {
		Convey("Should have customer record", func() {
			So(len(l), ShouldBeGreaterThan, 0)
		})
		Convey("Should be have 10 records", func() {
			So(len(l), ShouldEqual, 10)
		})
	})
}

func TestGetCustomerByID(t *testing.T) {
	removeAll()

	ID, v := buildRecord()
	err := models.AddCustomer(&v)
	value, errGetByID := models.GetCustomerByID(ID.Hex())

	Convey("Subject: Unit Test Get Customer by ID\n", t, func() {
		Convey("Should have created customer value", func() {
			So(err, ShouldEqual, nil)
		})
		Convey("Should have get customer by ID", func() {
			So(errGetByID, ShouldEqual, nil)
		})
		Convey("Should be same value initial and get by ID", func() {
			So(v, ShouldResemble, value)
		})
	})
}

func TestGetDeletedByID(t *testing.T) {
	removeAll()

	ID, v := buildRecord()
	errCreate := models.AddCustomer(&v)
	errDelete := models.DeleteCustomer(ID.Hex())
	value, err := models.GetCustomerByID(ID.Hex())

	Convey("Subject: Unit Test Get Customer by ID\n", t, func() {
		Convey("Should have created customer value", func() {
			So(errCreate, ShouldEqual, nil)
		})
		Convey("Should have deleted customer by ID", func() {
			So(errDelete, ShouldEqual, nil)
		})
		Convey("Should not have error when get record by ID", func() {
			So(err.Error(), ShouldEqual, "not found")
		})
		Convey("Should not have a record have beend deleted", func() {
			So(value, ShouldNotResemble, v)
		})
	})
}

func TestUpdateCustomer(t *testing.T) {
	removeAll()

	ID, v := buildRecord()

	u := structCustomer.Customer{
		Fullname: "ronyUpdated",
		Email:    "rony@fakerUpdated.com",
		Username: "ronyrusernameUpdated",
	}
	err := models.AddCustomer(&v)
	errUpdate := models.UpdateCustomer(ID.Hex(), &u)
	value, errGetByID := models.GetCustomerByID(ID.Hex())

	Convey("Subject: Unit Test Update Customer by ID\n", t, func() {
		Convey("Should have created customer value", func() {
			So(err, ShouldEqual, nil)
		})
		Convey("Should have update customer by ID", func() {
			So(errUpdate, ShouldEqual, nil)
		})
		Convey("Should have created or updated customer by ID", func() {
			So(errGetByID, ShouldEqual, nil)
		})
		Convey("Should have same ID, from initial ID and updated ID", func() {
			So(value.ID, ShouldEqual, ID)
		})
		Convey("Should have same value, from updated and the select updated",
			func() {
				So(value.Fullname, ShouldEqual, u.Fullname)
				So(value.Email, ShouldEqual, u.Email)
				So(value.Username, ShouldEqual, u.Username)
			})
	})
}
