package api

import "net/http"

// Endpoint is the starting point for all
// publishing activity
type Endpoint struct {
	location string
	password string
	client   *http.Client
}

// NewEndpoint starts a publishing session.
func NewEndpoint(authstring string, location string) (*Endpoint, error) {
	ep := &Endpoint{
		client:   &http.Client{},
		location: location + API_PREFIX,
		password: basicAuth(authstring),
	}

	return ep, nil
}
