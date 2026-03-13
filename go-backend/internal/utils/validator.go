package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// SetupValidator registers custom validators
func SetupValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("username", validateUsername)
		v.RegisterValidation("strongpassword", validateStrongPassword)
	}
}

// validateUsername checks: only letters, digits, underscore, dot; min 4 chars
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < 4 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_.]+$`, username)
	return matched
}

// validateStrongPassword checks: min 8, at least 1 upper, 1 lower, 1 digit
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	var hasUpper, hasLower, hasDigit bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		}
	}
	return hasUpper && hasLower && hasDigit
}

// ValidateUsername validates username manually and returns error message
func ValidateUsername(username string) error {
	if len(username) < 4 {
		return fmt.Errorf("username must be at least 4 characters")
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_.]+$`, username)
	if !matched {
		return fmt.Errorf("username can only contain letters, digits, underscores and dots")
	}
	return nil
}

// ValidatePassword validates password strength and returns error message
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	var hasUpper, hasLower, hasDigit bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		}
	}
	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	return nil
}

// FormatValidationErrors converts validator errors to readable map
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			field := toSnakeCase(e.Field())
			switch e.Tag() {
			case "required":
				errors[field] = fmt.Sprintf("%s is required", field)
			case "min":
				errors[field] = fmt.Sprintf("%s must be at least %s characters", field, e.Param())
			case "max":
				errors[field] = fmt.Sprintf("%s must be at most %s characters", field, e.Param())
			case "oneof":
				errors[field] = fmt.Sprintf("%s must be one of: %s", field, e.Param())
			default:
				errors[field] = fmt.Sprintf("%s is invalid", field)
			}
		}
	}

	return errors
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}
