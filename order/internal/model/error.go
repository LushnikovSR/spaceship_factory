package order

import "fmt"

type BaseError struct {
	// HTTP-code error.
	Code int `json:"code"`
	// Description of error.
	Message string `json:"message"`
}

func (e *BaseError) Error() string {
	return fmt.Sprintf("HTTP-code: %d Message: %s", e.Code, e.Message)
}

type NotFoundError struct {
	BaseError
}

type ConflictError struct {
	BaseError
}

type BadRequestError struct {
	BaseError
}

type InternalServerError struct {
	BaseError
}
