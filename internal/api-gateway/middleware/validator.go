package middleware

import (
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
)

type ValidatorManager struct {
	validators map[string]validator.Func
}

func NewValidatorManager() *ValidatorManager {
	mapValidators := make(map[string]validator.Func, 0)

	mapValidators["enum"] = checkValidEnum

	return &ValidatorManager{
		validators: mapValidators,
	}
}

func (vm *ValidatorManager) RegisterDefaultValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for tag, validatorFunc := range vm.validators {
			if err := v.RegisterValidation(tag, validatorFunc); err != nil {
				log.Fatalf("validator %s register error: %v", tag, err)
			}

			log.Printf("validator %s register success", tag)
		}
	}
}

func checkValidEnum(fl validator.FieldLevel) bool {
	enum, _ := fl.Field().Interface().(common.Enum)

	return enum.IsValid()
}
