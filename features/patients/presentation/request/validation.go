package request

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func ValidateBirthDate(fl validator.FieldLevel) bool {
	birthDate, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	_, err := time.Parse("2006-01-02", birthDate)
	return err == nil
}
