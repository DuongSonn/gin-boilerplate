package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	REQUIRED     = "required"
	EMAIL        = "email"
	PHONE_NUMBER = "phone_number"
	EQUAL_FIELD  = "eqfield"
)

var validatePhoneNumber validator.Func = func(fl validator.FieldLevel) bool {
	phoneNumber, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	return IsValidPhoneNumber(phoneNumber)
}

func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("phone_number", validatePhoneNumber)
}

func IsValidPhoneNumber(phoneNumber string) bool {
	// Check if the phone number matches the Vietnamese format
	vietnamesePhoneNumberPattern := `^(03[2-9]|07[0|6-9]|08[1-5]|09[0-9]|01[2|6|8|9])+([0-9]{8})$`
	match, err := regexp.MatchString(vietnamesePhoneNumberPattern, phoneNumber)
	if err != nil {
		return false
	}

	if !match {
		return false
	}

	return true
}

func IsValidEmail(email string) bool {
	// Check if the email matches the email format
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, err := regexp.MatchString(emailPattern, email)
	if err != nil {
		return false
	}

	if !match {
		return false
	}

	return true
}
