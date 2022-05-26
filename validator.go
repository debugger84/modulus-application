package application

import (
	"context"
	"errors"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type Validator interface {
	Validate(obj any) error
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
	return &DefaultValidator{validator: validate, translator: translator}
}

func (v *DefaultValidator) Validate(obj any) error {
	err := v.validator.Struct(obj)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Translate(v.translator))
		}
	}
	return nil
}
