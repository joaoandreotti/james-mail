package validator

import (
  "github.com/go-playground/validator/v10"
)

func ValidateStruct (data interface{}) bool {
  validate := validator.New()
  err := validate.Struct(data)
  return err == nil
}
