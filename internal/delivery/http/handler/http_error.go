package handler

import (
	"fmt"
	"strings"
)

type HTTPError struct {
	Err     error  `json:"-"`
	Code    int    `json:"-"`
	Message string `json:"error"`
} // @name HTTPError

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

func NewHTTPError(err error, code int, msg ...string) *HTTPError {
	msgs := strings.Join(msg, ", ")
	return &HTTPError{Err: err, Code: code, Message: msgs}
}
