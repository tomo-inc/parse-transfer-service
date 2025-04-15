package err

import "fmt"

var _ error = (*CustomError)(nil)

type CustomError struct {
	Status  int
	Code    int
	Message string
}

func NewCustomError(status int, code int, message string) *CustomError {
	return &CustomError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("status: %d code: %d, message: %s", e.Status, e.Code, e.Message)
}
