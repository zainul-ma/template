package models

import (
	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var session *mgo.Session

// Customer -- this user-service purpose for generic user
type (
	Customer struct {
		ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Fullname string        `json:"fullname,omitempty" bson:"fullname"`
		Username string        `json:"username,omitempt" bson:"username"`
		Email    string        `json:"email" bson:"email"`
	}
)

func getTableName() string {
	return "customer"
}

// Conn : initiate connection
func Conn() *mgo.Session {
	return session.Copy()
}

func init() {
	url := beego.AppConfig.String("mongodb:url")

	sess, err := mgo.Dial(url)

	if err != nil {

	}

	session = sess
	session.SetMode(mgo.Monotonic, true)
}

// GetCustomerByID method
func GetCustomerByID(id string) (cust Customer, err error) {
	connection := Conn()

	defer connection.Close()

	c := connection.DB("").C(getTableName())

	customer := Customer{}

	err = c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&customer)

	cust = customer

	return cust, err
}

// GetAllCustomer method
func GetAllCustomer(query map[string]string, fields []string, sortby []string,
	order []string, latestID string, limit int) (searchResults []interface{}, err error) {
	// fmt.Println("====")
	// fmt.Println(query)
	// fmt.Println("====")
	if latestID != "" {
		searchResults, err = Search(bson.M{
			"_id": bson.M{
				"&gt": bson.ObjectIdHex(latestID),
			},
		}, limit)
	} else {
		searchResults, err = Search(nil, limit)
	}

	return searchResults, err
}

// withCollection method
func withCollection(s func(*mgo.Collection) error) error {
	connection := Conn()

	defer connection.Close()

	c := connection.DB("").C(getTableName())
	return s(c)
}

//Search for base method
func Search(q interface{}, limit int) (searchResults []interface{},
	searchErr error) {
	query := func(c *mgo.Collection) error {
		fn := c.Find(q).Limit(limit).All(&searchResults)
		if limit < 0 {
			fn = c.Find(q).All(&searchResults)
		}
		return fn
	}
	search := func() error {
		return withCollection(query)
	}
	err := search()
	if err != nil {
		searchErr = err
	}
	return
}

//UpdateCustomer method
func UpdateCustomer(id string, customer *Customer) (err error) {
	connection := Conn()

	defer connection.Close()

	c := connection.DB("").C(getTableName())

	err = c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, customer)

	return
}

// AddCustomer method
func AddCustomer(customer *Customer) (err error) {
	connection := Conn()

	defer connection.Close()

	c := connection.DB("").C(getTableName())

	err = c.Insert(customer)

	return
}

// DeleteCustomer method
func DeleteCustomer(id string) (err error) {
	connection := Conn()

	defer connection.Close()

	c := connection.DB("").C(getTableName())

	err = c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})

	return
}
