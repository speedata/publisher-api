package api

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	apiPrefix string = "/v0"
)

// Endpoint is the starting point for all
// publishing activity
type Endpoint struct {
	location string
	password string
	client   *http.Client
}

// AvailableVersions returns a list of versions available on the server. The
// version 'latest' is not included in the list.
func (ep *Endpoint) AvailableVersions() ([]string, error) {
	loc := ep.location + "/versions"
	req, err := http.NewRequest("GET", loc, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(ep.password, "")
	resp, err := ep.client.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var versions []string
	err = json.Unmarshal(data, &versions)
	if err != nil {
		return nil, err
	}

	return versions, nil
}

// NewEndpoint starts a publishing session.
func NewEndpoint(authstring string, location string) (*Endpoint, error) {
	ep := &Endpoint{
		client:   &http.Client{},
		location: location + apiPrefix,
		password: authstring,
	}

	return ep, nil
}
