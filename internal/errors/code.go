package errors

import (
	"fmt"
	"net/http"
)

const (
	dynamodbVersionPrefix = "com.amazonaws.dynamodb.v20120810"
)

func WithBadRequest(msg string) WrapOption {
	return func(e *Error) {
		e.Type = "BADREQUEST"
		e.Message = msg
		e.StatusCode = http.StatusBadRequest
	}
}

func WithConditionalCheckFailedException() WrapOption {
	return func(e *Error) {
		e.Type = fmt.Sprintf("%s#%s", dynamodbVersionPrefix, "ConditionalCheckFailedException")
		e.Message = "Failed condition."
		e.StatusCode = http.StatusBadRequest
	}
}

// WithCodeUnknown is a generic error that can be used as a last
// resort if there is no situation-specific error message that can be used
func WithCodeUnknown() WrapOption {
	return func(e *Error) {
		e.Type = "UNKNOWN"
		e.Message = "unknown error"
		e.StatusCode = http.StatusInternalServerError
	}
}

func WithNotFound() WrapOption {
	return func(e *Error) {
		e.Type = "NOTFOUND"
		e.Message = "not found"
		e.StatusCode = http.StatusNotFound
	}
}


func WithInternalServerError() WrapOption {
	return func(e *Error) {
		e.Type = "NOTFOUND"
		e.Message = "internal server error"
		e.StatusCode = http.StatusInternalServerError
	}
}
