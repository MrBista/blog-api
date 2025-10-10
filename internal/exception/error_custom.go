package exception

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ValidationErrorItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorCustom struct {
	Message string                `json:"message"`
	Errors  []ValidationErrorItem `json:"errors,omitempty"`
	Status  int                   `json:"status"`
	Code    string                `json:"code"`
}

func (e *ErrorCustom) Error() string {
	return e.Message
}

func (e *ErrorCustom) GetStatusHttp() int {
	if e.Status == 0 {
		return fiber.StatusInternalServerError
	}

	return e.Status
}

const (
	ERR_BAD_REQUEST   = "BAD_REQUEST_ERROR"
	ERR_VALIDATION    = "VALIDATION_ERROR"
	ERR_BUSNISS_LOGIC = "BUSNIS_LOGIC_ERROR"
	ERR_DB            = "DB_ERROR"
	ERR_UNHANDLE      = "INTERNAL_SERVER_ERROR"
	ERR_COMMON        = "COMMON_ERR"
)

func NewBadRequestErr(message string) *ErrorCustom {
	return &ErrorCustom{
		Message: message,
		Status:  fiber.StatusBadRequest,
		Code:    ERR_BAD_REQUEST,
	}
}

func NewValidationErr(err error) *ErrorCustom {
	var errors []ValidationErrorItem

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, filedErr := range validationErrs {
			errors = append(errors, ValidationErrorItem{
				Field:   filedErr.Field(),
				Message: filedErr.Error(),
			})
		}
	}

	return &ErrorCustom{
		Errors:  errors,
		Message: "Validation error, make sure body is correct",
		Status:  fiber.StatusBadRequest,
		Code:    ERR_BAD_REQUEST,
	}
}
