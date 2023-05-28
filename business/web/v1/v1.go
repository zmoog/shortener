package v1

import "errors"

type ErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

type RequestError struct {
	Err    error
	Status int
}

func NewRequestError(err error, status int) error {
	return &RequestError{Err: err, Status: status}
}

func (r *RequestError) Error() string {
	return r.Err.Error()
}

func GetRequestError(err error) *RequestError {
	var re *RequestError
	if !errors.As(err, &re) {
		return nil
	}
	return re
}

func IsRequestError(err error) bool {
	var re *RequestError
	return errors.As(err, &re)
}
