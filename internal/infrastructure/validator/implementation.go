package validator

import (
	"net/url"

	govalidator "github.com/go-playground/validator/v10"
)

type validate struct {
	validator *govalidator.Validate
}

func NewValidator() *validate {
	v := validate{
		validator: govalidator.New(),
	}
	if err := v.validator.RegisterValidation(
		"url_without_scheme", validateURLWithoutScheme); err != nil {
		panic(err)
	}

	return &v
}

func (va *validate) ValidateStruct(v any) error {
	return va.validator.Struct(v)
}

func validateURLWithoutScheme(fl govalidator.FieldLevel) bool {
	u := fl.Field().String()
	if u == "" {
		return false
	}

	urlParsed, err := url.Parse(u)
	if err != nil {
		return false
	}
	if urlParsed.Scheme == "" {
		u = "http://" + u
	}

	_, err = url.ParseRequestURI(u)

	return err == nil
}
