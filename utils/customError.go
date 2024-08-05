package utils

import "fmt"

type RequestError struct {
	StatusCode int

	Err error
}

func (r *RequestError) Error() string {
	return fmt.Sprintf(" %v", r.Err)
}
