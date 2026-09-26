package restful

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	globalValidator *validator.Validate
	validatorOnce   sync.Once
)

func GetValidator() *validator.Validate {
	validatorOnce.Do(func() {
		globalValidator = validator.New()
	})
	return globalValidator
}
func SetValidator(v *validator.Validate) {
	globalValidator = v
}

type ValidateError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   any    `json:"value,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *ValidateError) Error() string {
	return fmt.Sprintf("field '%s' failed on '%s'", e.Field, e.Tag)
}

type ValidateErrors []*ValidateError

func (e ValidateErrors) Error() string {
	msgs := make([]string, len(e))
	for i, ve := range e {
		msgs[i] = ve.Error()
	}
	return strings.Join(msgs, "; ")
}

func (r *Request) BindAndValidate(v any) error {
	if err := r.ReadEntity(v); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}
	return ValidateStruct(v)
}

func ValidateStruct(v any) error {
	if v == nil {
		return errors.New("validation target is nil")
	}

	if validatable, ok := v.(Validatable); ok {
		return validatable.Validate()
	}

	return validateWithPlayground(v)
}

func validateWithPlayground(v any) error {
	err := GetValidator().Struct(v)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err
	}

	out := make(ValidateErrors, 0, len(validationErrors))
	for _, fieldErr := range validationErrors {
		out = append(out, &ValidateError{
			Field:   fieldErr.Field(),
			Tag:     fieldErr.Tag(),
			Value:   fieldErr.Value(),
			Message: formatFieldError(fieldErr),
		})
	}
	return out
}

func formatFieldError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", fe.Field(), fe.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s", fe.Field(), fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be >= %s", fe.Field(), fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be <= %s", fe.Field(), fe.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", fe.Field(), fe.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", fe.Field())
	default:
		return fmt.Sprintf("%s failed on '%s'", fe.Field(), fe.Tag())
	}
}
