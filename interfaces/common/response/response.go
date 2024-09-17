package response

import (
	"github.com/google/uuid"
)

func GetUUIDParam(param string) (uuid.UUID, error) {
	return uuid.Parse(param)
}

type ApiResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Errors  []string    `json:"errors"`
	Data    interface{} `json:"data"`
}

func BuildResponse(code string, message string, errors []string, data interface{}) *ApiResponse {
	r := &ApiResponse{}
	r.Code = code
	r.Message = message
	r.Errors = errors
	r.Data = data
	return r
}

type FieldValidationError struct {
	Namespace string      `json:"namespace"`
	Field     string      `json:"field"`
	Error     string      `json:"error"`
	Kind      string      `json:"kind"`
	Value     interface{} `json:"value"`
}
