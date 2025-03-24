package utils

import "fmt"

type CustomError struct {
	Message string
	Code    int
}

func (e CustomError) Error() string {
	return fmt.Sprintf("message: %s\n", e.Message)
}
