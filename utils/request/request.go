package request

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

func ValidateRequest(err error) any {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make(map[string]string)
		for _, fe := range ve {
			out[fe.Field()] = "This field is required"
		}

		return out
	} else {
		return "invalid payload"
	}
}
