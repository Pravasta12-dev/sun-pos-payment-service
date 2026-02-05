package validator

import (
	"errors"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
)

type Validator struct {
	Validator  *validator.Validate
	Translator ut.Translator
}

func NewValidator() *Validator {
	en := en.New()
	uni := ut.New(en, en)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatalf("Translator Not Found")
	}

	validate := validator.New()

	return &Validator{
		Validator:  validate,
		Translator: trans,
	}
}

func (v *Validator) Validate(s interface{}) error {
	err := v.Validator.Struct(s)

	if err != nil {
		object, _ := err.(validator.ValidationErrors)
		for _, e := range object {
			log.Infof("[Validator-1] field %s, tag %s, actualTag %s, param %s, value %v",
				e.Field(),
				e.Tag(),
				e.ActualTag(),
				e.Param(),
				e.Value(),
			)

			return errors.New(e.Translate(v.Translator))
		}
	}

	return nil
}
