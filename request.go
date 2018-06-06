package main

import (
	"encoding/json"
	"io"

	"github.com/asaskevich/govalidator"
)

type RequestSignUp struct {
	Email     string `json:"email" valid:"email,required"`
	Password  string `json:"password" valid:"required"`
	FirstName string `json:"first_name" valid:"required"`
	LastName  string `json:"last_name" valid:"required"`
}

// New request sign up from io reader
func NewRequestSignUp(r io.Reader) (*RequestSignUp, error) {
	data := new(RequestSignUp)
	decoder := json.NewDecoder(r)
	err := decoder.Decode(data)
	return data, err
}

// Validate SignUp request, and return error messages if not valid
func (d *RequestSignUp) Validate() (bool, map[string]string) {
	isValid, err := govalidator.ValidateStruct(d)
	errs := govalidator.ErrorsByField(err)

	return isValid, errs
}
