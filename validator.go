package application

import (
	"context"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"reflect"
	"strings"
)

type Validator interface {
	Validate(obj any) []ValidationError
}

type DefaultValidator struct {
	validator  *validator.Validate
	translator ut.Translator
	logger     Logger
}

func NewDefaultValidator(logger Logger) Validator {
	uni := ut.New(en.New())
	translator, _ := uni.GetTranslator("en")
	validate := validator.New()
	err := enTranslations.RegisterDefaultTranslations(validate, translator)
	if err != nil {
		logger.Error(context.Background(), "Cannot register default translations for validator")
	}
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &DefaultValidator{validator: validate, translator: translator}
}

func (v *DefaultValidator) Validate(obj any) []ValidationError {
	err := v.validator.Struct(obj)
	if err != nil {
		if validatorErr, ok := err.(validator.ValidationErrors); ok {
			result := make([]ValidationError, len(validatorErr))
			for i, validationError := range validatorErr {
				result[i] = *NewValidationError(
					validationError.Field(),
					validationError.Translate(v.translator),
				)
			}
			return result
		} else {
			return []ValidationError{*NewValidationError(
				"",
				err.Error(),
			)}
		}
	}
	return nil
}
