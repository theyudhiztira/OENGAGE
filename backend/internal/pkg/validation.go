package pkg

import "github.com/go-playground/validator/v10"

func ParseValidationMessage(err error) []string {
	var messages []string
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrors {
				messages = append(messages, "Field '"+fieldErr.Field()+"' failed on the '"+fieldErr.Tag()+"' rule")
			}
		}

	}
	return messages
}
