package web

import "errors"

type shutdownError struct {
	Message string
}

func NewShutdownError(message string) error {
	return &shutdownError{Message: message}
}

func (s *shutdownError) Error() string {
	return s.Message
}

func IsShutdownError(err error) bool {
	var se *shutdownError
	return errors.As(err, &se)
}
