package util

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func ValidateISO8601(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	_, err := time.Parse(DateFormat, dateStr)
	return err == nil
}
