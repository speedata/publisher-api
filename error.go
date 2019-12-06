package api

import (
	"fmt"
)

// Error is an error struct which is used when the communication between
// the client and the endpoint is broken.
type Error struct {
	ErrorType string `json:"type"`
	Title     string
	Detail    string
	Instance  string
	RequestID int
}

// Give a meaningful string of the error message.
func (a Error) Error() string {
	return fmt.Sprintf("Publishing error: %s", a.Title)
}
