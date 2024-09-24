package errorx

import (
	"fmt"
	"net/http"
)

type TektonError struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (e *TektonError) Error() string {
	return e.Message
}

func NewError(code int64, message string, data interface{}) error {
	return &TektonError{Code: code, Message: message, Data: data}
}

func NewDefaultError(message string, a ...any) error {
	return &TektonError{Code: http.StatusInternalServerError, Message: fmt.Sprintf(message, a...)}
}
