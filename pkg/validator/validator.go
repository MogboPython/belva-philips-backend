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

const splitInt = 2

// New creates a new validator
func New() *Validator {
	validate := validator.New()

	// Register validation for struct fields
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", splitInt)[0]
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
func (v *Validator) Validate(i any) error {
	if err := v.validate.Struct(i); err != nil {
		var errors []string

		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		for _, err := range validationErrors {
			var message string

			switch err.Tag() {
			case "required":
				message = err.Field() + " is required"
			case "email":
				message = err.Field() + " must be a valid email"
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
