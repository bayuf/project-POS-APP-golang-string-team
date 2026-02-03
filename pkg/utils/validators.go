package utils

import "github.com/go-playground/validator/v10"

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateErrors(data any) ([]FieldError, error) {
	validate := validator.New()

	err := validate.Struct(data)
	if err == nil {
		return nil, nil
	}

	var errors []FieldError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {

			field := err.Field()
			if jsonTag := err.StructField(); jsonTag != "" {
				field = err.Field()
			}

			var message string
			switch err.Tag() {
			case "required":
				message = field + " is required"
			case "min":
				message = field + " must be at least " + err.Param() + " characters"
			case "max":
				message = field + " must be at most " + err.Param() + " characters"
			case "email":
				message = field + " must be a valid email address"
			case "alphanum":
				message = field + " must contain only alphanumeric characters"
			default:
				message = field + " is invalid"
			}

			errors = append(errors, FieldError{
				Field:   field,
				Message: message,
			})
		}
		return errors, err
	}

	return nil, err
}
