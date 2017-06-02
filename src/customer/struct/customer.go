package structcustomer

import (
  "gopkg.in/go-playground/validator.v9"
  "fmt"
)

type TypeCustomerRequest struct {
  Id int `validate:"required"`
  Fullname string `validate:"required"`
  Username string `validate:"required"`
  Email string `validate:"required"`
}

var validate *validator.Validate

func ValidateTypeCustomerRequest(TypeCustomerRequest *TypeCustomerRequest) {
  err := validate.Struct(TypeCustomerRequest)
}

func CheckErrValidate(err error){
  if err != nil {
    if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace()) // can differ when a custom TagNameFunc is registered or
			fmt.Println(err.StructField())     // by passing alt name to ReportError like below
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}

		// from here you can create your own error messages in whatever language you wish
    return
  }
}
