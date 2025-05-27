package utils

import (
	"errors"
	"fmt"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type ApiError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func castValidationError(arrayErrors validator.ValidationErrors) []ApiError {
	res := make([]ApiError, len(arrayErrors))

	for idx, fe := range arrayErrors {
		res[idx] = ApiError{Field: fe.Field(), Message: msgForTag(fe)}
	}

	return res
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("The '%s' field is required.", fe.Field())
	case "email":
		return fmt.Sprintf("The '%s' field must be a valid email address.", fe.Field())
	case "oneof":
		return fmt.Sprintf("The value of '%s' must be one of the following: [%s].", fe.Field(), strings.Join(strings.Split(fe.Param(), " "), ", "))
	case "numeric":
		return fmt.Sprintf("The '%s' field must be a number.", fe.Field())
	case "alphanum":
		return fmt.Sprintf("The '%s' field can only contain letters and numbers.", fe.Field())
	case "gte":
		return fmt.Sprintf("The '%s' field must be greater than or equal to '%v'.", fe.Field(), fe.Param())
	case "min":
		return fmt.Sprintf("The '%s' field must be at least %v characters long.", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("The '%s' field cannot exceed %v characters.", fe.Field(), fe.Param())
	case "alpha":
		return fmt.Sprintf("The '%s' field can only contain letters.", fe.Field())
	case "len":
		return fmt.Sprintf("The '%s' field must be equal %v characters.", fe.Field(), fe.Param())
	case "uuid":
		return fmt.Sprintf("The '%s' field must be a valid uuid.", fe.Field())
	case "enum":
		enum, _ := fe.Value().(common.Enum)
		return enum.ErrorMessage()
	default:
		return fmt.Sprintf("The '%s' field is invalid.", fe.Field())
	}
}

type TechnicalError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

type BusinessError struct {
	Code      int    `json:"-"`
	Message   string `json:"message"`
	ErrorCode string `json:"error_code"`
}

func (err TechnicalError) Error() string {
	return err.Message
}

func (err BusinessError) Error() string {
	return err.Message
}

func HandleValidateData(ctx *gin.Context, err error) {
	var targetError validator.ValidationErrors

	if errors.As(err, &targetError) {
		apiErrors := castValidationError(targetError)
		ErrorResponse(ctx, http.StatusBadRequest, apiErrors)
		return
	}

	ErrorResponse(ctx, http.StatusBadRequest, ApiError{
		Field:   "",
		Message: err.Error(),
	})
}

func HandleErrorResponse(ctx *gin.Context, err error) {
	var techError TechnicalError
	var businessError BusinessError

	if errors.As(err, &techError) {
		ErrorResponse(ctx, techError.Code, techError)
		return
	}

	if errors.As(err, &businessError) {
		ErrorResponse(ctx, businessError.Code, businessError)
		return
	}

	ErrorResponse(ctx, http.StatusInternalServerError, common.MSG_INTERNAL_ERROR)
}
