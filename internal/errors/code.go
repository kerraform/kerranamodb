package errors

import "net/http"

func WithBadRequest() WrapOption {
	return func(e *Error) {
		e.Code = "BADREQUEST"
		e.Message = "request invalid"
		e.StatusCode = http.StatusBadRequest
	}
}

// WithCodeUnknown is a generic error that can be used as a last
// resort if there is no situation-specific error message that can be used
func WithCodeUnknown() WrapOption {
	return func(e *Error) {
		e.Code = "UNKNOWN"
		e.Message = "unknown error"
		e.StatusCode = http.StatusInternalServerError
	}
}

func WithNotFound() WrapOption {
	return func(e *Error) {
		e.Code = "NOTFOUND"
		e.Message = "not found"
		e.StatusCode = http.StatusNotFound
	}
}
