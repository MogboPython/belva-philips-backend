package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the validator.Validate struct
type Validator struct {
	validate *validator.Validate
}

// New creates a new validator
func New() *Validator {
	validate := validator.New()

	// Register validation for struct fields
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{
		validate: validate,
	}
}

// Validate validates the given struct
func (v *Validator) Validate(i interface{}) error {
	if err := v.validate.Struct(i); err != nil {
		var errors []string

		for _, err := range err.(validator.ValidationErrors) {
			var message string

			switch err.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", err.Field())
			case "email":
				message = fmt.Sprintf("%s must be a valid email", err.Field())
			case "min":
				message = fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
			case "max":
				message = fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param())
			default:
				message = fmt.Sprintf("%s failed validation: %s", err.Field(), err.Tag())
			}

			errors = append(errors, message)
		}

		return fmt.Errorf("validation failed: %s", strings.Join(errors, ", "))
	}

	return nil
}
