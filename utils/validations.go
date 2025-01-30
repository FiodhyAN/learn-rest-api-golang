package utils

import (
	"unicode"

	"github.com/FiodhyAN/learn-rest-api-golang/types"
	"github.com/go-playground/validator/v10"
)

func validatePassword(fl validator.FieldLevel) bool {
	payload := fl.Parent().Interface().(types.RegisterPayload)
	var (
		hasUpperCase bool
		hasLowerCase bool
		hasDigit     bool
		hasSpecial   bool
	)

	for _, char := range payload.Password {
		if unicode.IsUpper(char) {
			hasUpperCase = true
		} else if unicode.IsLower(char) {
			hasLowerCase = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else {
			hasSpecial = true
		}
	}

	if hasUpperCase && hasLowerCase && hasDigit && hasSpecial {
		return true
	} else {
		return false
	}
}
