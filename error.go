package api

import (
	"fmt"
)

type APIError struct {
	ErrorType string `json:"type"`
	Title     string
	Detail    string
	Instance  string
}

func (a APIError) Error() string {
	return fmt.Sprintf("Publishing error")
}
