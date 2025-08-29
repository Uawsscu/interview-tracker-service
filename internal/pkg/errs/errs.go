package errs

import "net/http"

type HttpError struct {
	Code    int
	Message string
}

func (e *HttpError) Error() string {
	return e.Message
}

func Unauthorized(msg string) *HttpError {
	return &HttpError{Code: http.StatusUnauthorized, Message: msg}
}
