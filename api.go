package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

const (
	API_PREFIX string = "/v0"
)

func basicAuth(username string) string {
	auth := username + ":"
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// Endpoint is the starting point for all
// publishing activity
type Endpoint struct {
	location string
	password string
	client   *http.Client
}

// PublishRequest is
type PublishRequest struct {
	endpoint *Endpoint
	Files    []PublishFile
}

// AttachFile adds a file to the PublishRequest
func (p *PublishRequest) AttachFile(pathToFile string) error {
	filename := filepath.Base(pathToFile)
	data, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return err
	}

	p.Files = append(p.Files, PublishFile{Filename: filename, Contents: data})
	return nil
}

// PublishFile is a file for the publishing request
type PublishFile struct {
	Filename string
	Contents []byte
}

// PublishResponse holds the id to the publishing process
type PublishResponse struct {
	endpoint *Endpoint
	Id       string
}

// GetPDF gets the PDF from the server. In case of an error, the byte slice might not be meaningful.
// Otherwise it holds the PDF file.
func (p *PublishResponse) GetPDF() ([]byte, error) {
	loc := p.endpoint.location + "/pdf/" + p.Id
	req, err := http.NewRequest("GET", loc, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+p.endpoint.password)

	resp, err := p.endpoint.client.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
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

// Publish sends data to the server.
func (e *Endpoint) Publish(data *PublishRequest) (PublishResponse, error) {
	var p PublishResponse

	loc := e.location + "/publish"

	b, err := json.Marshal(data)
	if err != nil {
		return p, err
	}

	br := bytes.NewReader(b)

	req, err := http.NewRequest("POST", loc, br)
	if err != nil {
		return p, err
	}
	req.Header.Add("Authorization", "Basic "+e.password)
	req.Header.Add("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return p, err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return p, err
	}

	err = json.Unmarshal(buf, &p)
	if err != nil {
		return p, err
	}
	p.endpoint = e

	return p, nil
}

// NewPublishRequest is the base structure to start a publishing request to the endpoint.
func (e *Endpoint) NewPublishRequest() *PublishRequest {
	p := &PublishRequest{
		endpoint: e,
	}
	return p
}
