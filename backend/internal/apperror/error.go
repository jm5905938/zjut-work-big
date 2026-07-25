package apperror

import "net/http"

type Error struct {
	Status  int
	Code    int
	Message string
	Err     error
}

func New(status int, code int, message string, err error) *Error {
	return &Error{
		Status:  status,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func BadRequest(message string) *Error {
	return New(http.StatusBadRequest, 400, message, nil)
}

func Unauthorized(message string) *Error {
	return New(http.StatusUnauthorized, 401, message, nil)
}

func Forbidden(message string) *Error {
	return New(http.StatusForbidden, 403, message, nil)
}

func NotFound(message string) *Error {
	return New(http.StatusNotFound, 404, message, nil)
}

func Conflict(message string) *Error {
	return New(http.StatusConflict, 409, message, nil)
}

func Internal(err error) *Error {
	return New(
		http.StatusInternalServerError,
		500,
		"服务器死了啊）",
		err,
	)
}
