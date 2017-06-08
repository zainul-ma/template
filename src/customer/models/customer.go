package models

import (
	"errors"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	structCustomer "customer/structs"
)

// GetTableName method
func GetTableName() string {
	return "customer"
}

// GetDBName method
func GetDBName() string {
	return "customer"
}

// GetCustomerByID method
func GetCustomerByID(id string) (cust structCustomer.Customer, err error) {
	connection := ConnectMongo()

	defer connection.Close()

	c := connection.DB(GetDBName()).C(GetTableName())

	customer := structCustomer.Customer{}

	err = c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&customer)

	cust = customer

	return cust, err
}

// GetAllCustomer method
func GetAllCustomer(query map[string]string, fields []string, sortby []string,
	order []string, latestID string,
	limit int) (searchResults []interface{}, err error) {

	queryBuilder := bson.M{}
	if len(query) != 0 {
		queryBuilder["$or"] = []bson.M{}
	}

	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		queryBuilder["$or"] = append(queryBuilder["$or"].([]bson.M),
			bson.M{k: v})
	}

	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New(
						"Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}

		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order,
			// all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New(
						"Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New(
				"Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	if latestID != "" {
		queryBuilder["_id"] = bson.M{
			"&gt": bson.ObjectIdHex(latestID),
		}

		searchResults, err = Search(queryBuilder, fields, sortFields, limit)
	} else {
		searchResults, err = Search(queryBuilder, fields, sortFields, limit)
	}

	return searchResults, err
}

// withCollection method
func withCollection(s func(*mgo.Collection) error) error {
	connection := ConnectMongo()
	defer connection.Close()
	c := connection.DB(GetDBName()).C(GetTableName())
	return s(c)
}

// selectOnlyFields method
func selectOnlyFields(q ...string) (r bson.M) {
	r = make(bson.M, len(q))
	for _, s := range q {
		r[s] = 1
	}
	return
}

//Search for base method
func Search(q interface{}, fields []string, sortBy []string, limit int) (
	searchResults []interface{}, searchErr error) {
	query := func(c *mgo.Collection) error {
		fn := c.Find(q).Select(
			selectOnlyFields(fields...)).Sort(
			sortBy...).Limit(limit).All(&searchResults)

		if limit < 0 {
			fn = c.Find(q).Select(
				selectOnlyFields(fields...)).Sort(
				sortBy...).All(&searchResults)
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
func UpdateCustomer(id string, customer *structCustomer.Customer) (err error) {
	connection := ConnectMongo()
	defer connection.Close()
	c := connection.DB(GetDBName()).C(GetTableName())
	err = c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, customer)
	return
}

// AddCustomer method
func AddCustomer(customer *structCustomer.Customer) (err error) {
	connection := ConnectMongo()
	defer connection.Close()
	c := connection.DB(GetDBName()).C(GetTableName())
	err = c.Insert(customer)
	return
}

// DeleteCustomer method
func DeleteCustomer(id string) (err error) {
	connection := ConnectMongo()
	defer connection.Close()
	c := connection.DB(GetDBName()).C(GetTableName())
	err = c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	return
}
