package internalerrors

import (
	"errors"

	"github.com/go-playground/validator"
)

func ValidateStruct(obj interface{}) error {
	validate := validator.New()
	err := validate.Struct(obj)

	if err == nil {
		return nil
	}

	validationErrors := err.(validator.ValidationErrors)
	validationError := validationErrors[0]

	switch validationError.Tag() {
	case "required":
		return errors.New(validationError.StructField() + " is required")
	case "min":
		return errors.New(validationError.StructField() + " is required with min " + validationError.Param())
	case "max":
		return errors.New(validationError.StructField() + " is required with max " + validationError.Param())
	case "email":
		return errors.New(validationError.StructField() + " is invalid")
	}
	return nil
}
