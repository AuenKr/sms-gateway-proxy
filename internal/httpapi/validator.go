package httpapi

import (
	"github.com/go-playground/validator/v10"
	validators "github.com/go-playground/validator/v10/non-standard/validators"
	"go.uber.org/fx"
)

type NewValidatorResult struct {
	fx.Out

	Validate *validator.Validate
}

func NewValidator() NewValidatorResult {
	validate := validator.New()

	validate.RegisterValidation("notblank", validators.NotBlank)

	return NewValidatorResult{Validate: validate}
}
