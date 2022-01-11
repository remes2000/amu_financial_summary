package validators

import (
	"github.com/go-playground/validator/v10"
	"time"
)

var ValidDateLayout = "02-01-2006"

var ValidDate validator.Func = func(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(string)
	_, timeParseErr := time.Parse(ValidDateLayout, value)
	return timeParseErr == nil
}
