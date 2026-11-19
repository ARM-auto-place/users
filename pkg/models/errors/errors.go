package errors

import "net/http"

type Error struct {
	status int
	code   string
	reason string
	cause  error
}

func (e *Error) Error() string         { return e.reason }
func (e *Error) Unwrap() error         { return e.cause }
func (e *Error) GetStatusCode() int    { return e.status }
func (e *Error) GetPublicText() string { return "users.errors." + e.code }

func BadRequest(reason string) *Error {
	return &Error{status: http.StatusBadRequest, code: "badRequest", reason: reason}
}

func Unauthorized() *Error {
	return &Error{status: http.StatusUnauthorized, code: "invalidSession", reason: "a valid user session is required"}
}

func Forbidden() *Error {
	return &Error{status: http.StatusForbidden, code: "forbidden", reason: "origin is not allowed"}
}

func ProfileNotFound() *Error {
	return &Error{status: http.StatusNotFound, code: "profileNotFound", reason: "profile not found"}
}

func InternalServerError(cause error) *Error {
	return &Error{status: http.StatusInternalServerError, code: "internalServerError", reason: "internal server error", cause: cause}
}
