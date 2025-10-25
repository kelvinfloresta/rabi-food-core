package utils

import (
	"fmt"

	ptbr_translations "github.com/go-playground/validator/v10/translations/pt_BR"

	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	Validator  = validator.New()
	translator ut.Translator
)

func init() {
	locale := pt_BR.New()
	uni := ut.New(locale, locale)
	trans, _ := uni.GetTranslator("en")
	translator = trans

	err := ptbr_translations.RegisterDefaultTranslations(Validator, translator)
	if err != nil {
		panic(fmt.Sprintf("Failed to register translations: %v", err))
	}
}

func TranslateValidationErrors(err validator.ValidationErrors) []string {
	errors := make([]string, 0, len(err))
	for _, e := range err {
		errors = append(errors, e.Translate(translator))
	}

	return errors
}
