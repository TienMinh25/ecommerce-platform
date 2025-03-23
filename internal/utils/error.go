package utils

import (
	"errors"
	"fmt"
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
		return fmt.Sprintf("Trường '%s' là bắt buộc.", fe.Field())
	case "email":
		return fmt.Sprintf("Trường '%s' phải là một địa chỉ email hợp lệ.", fe.Field())
	case "oneof":
		return fmt.Sprintf("Giá trị của trường '%s' phải thuộc một trong các lựa chọn: [%s].", fe.Field(), strings.Join(strings.Split(fe.Param(), " "), ", "))
	case "numeric":
		return fmt.Sprintf("Trường '%s' phải là số.", fe.Field())
	case "alphanum":
		return fmt.Sprintf("Trường '%s' chỉ được chứa chữ cái và số.", fe.Field())
	case "gte":
		return fmt.Sprintf("Trường '%s' phải lớn hơn hoặc bằng '%v'", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("Trường '%s' không hợp lệ.", fe.Field())
	}
}

type TechnicalError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BusinessError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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
		ErrorResponse(ctx, techError.Code, techError.Message)
		return
	}

	if errors.As(err, &businessError) {
		ErrorResponse(ctx, businessError.Code, businessError.Message)
	}
}
