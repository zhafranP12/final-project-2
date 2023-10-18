package helpers

import (
	"finalProject2/pkg/errs"
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

func ValidateStruct(payload interface{}) errs.Error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	// Translator for validation error
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)

	fmt.Println(payload)
	err := validate.Struct(payload)
	fmt.Println(err)

	if err != nil {
		var msg string
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			msg = fieldError.Translate(trans)
			return errs.NewBadRequest(msg)
		}
	}
	return nil
}
