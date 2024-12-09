package validator

import (
	"errors"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

//nolint:gochecknoinits
func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	english := en.New()
	uni := ut.New(english, english)
	trans, _ = uni.GetTranslator("en")
	if err := enTranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		panic(err)
	}
}

func Var(v any, tag string) error {
	return buildError(validate.Var(v, tag))
}

func Struct(v any) error {
	return buildError(validate.Struct(v))
}

func buildError(err error) error {
	if nil == err {
		return nil
	}

	if ve, ok := err.(validator.ValidationErrors); ok { //nolint:errorlint
		errMsgs := make([]string, len(ve))
		for i, fe := range ve {
			errMsgs[i] = fe.Translate(trans)
		}
		return errors.New(strings.Join(errMsgs, ", "))
	}

	return err
}
