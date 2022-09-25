package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

var (
	// ErrNotFound is returned on a 404
	ErrNotFound error = fmt.Errorf("Resource not found")
)

func basicAuth(username string) string {
	auth := username + ":"
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// PublishRequest is an instance to send data to a server and get a PDF. One
// Endpoint can have multiple PublishingRequests.
type PublishRequest struct {
	endpoint *Endpoint
	Version  string
	Files    []PublishFile `json:"files"`
}

// AttachFile adds a file to the PublishRequest for the server. Usually you have
// to provide the layout and the data file. All assets (fonts, images) can be
// referenced by http(s) hyperlinks.
func (p *PublishRequest) AttachFile(pathToFile string) error {
	filename := filepath.Base(pathToFile)
	data, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return err
	}

	p.Files = append(p.Files, PublishFile{Filename: filename, Contents: data})
	return nil
}

// PublishFile is a file for the publishing request.
type PublishFile struct {
	Filename string `json:"filename"`
	Contents []byte `json:"contents"`
}

// PublishResponse holds the id to the publishing process.
type PublishResponse struct {
	endpoint *Endpoint
	ID       string
}

// Errormessage contains a message from the publisher together with its error
// code. The error message is a message from the publishing run (like image not
// found).
type Errormessage struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

// ProcessStatus contains information about the current status of the PDF
// generation. If the Finished field is nil, Errors and Errormessages are not
// set.
type ProcessStatus struct {
	Finished      *time.Time
	Errors        int
	Errormessages []Errormessage
}

// Status returns the status of the publishing run. If the process is still
// running, the Finished field is set to nil.
func (p *PublishResponse) Status() (*ProcessStatus, error) {
	loc := p.endpoint.location + "/status/" + p.ID
	req, err := http.NewRequest("GET", loc, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(p.endpoint.password, "")

	resp, err := p.endpoint.client.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ps := &ProcessStatus{}
	err = json.Unmarshal(buf, ps)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

// Wait for the publishing process to finish. Return an error if something is
// wrong with the request. If there is an error during the publishing run but
// the request itself is without errors, the error is nil, but the returned
// publishing status has the numbers of errors set.
func (p *PublishResponse) Wait() (*ProcessStatus, error) {
	loc := p.endpoint.location + "/wait/" + p.ID
	req, err := http.NewRequest("GET", loc, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(p.endpoint.password, "")

	resp, err := p.endpoint.client.Do(req)
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 422:
		var ae Error
		err = json.Unmarshal(buf, &ae)
		if err != nil {
			return nil, err
		}
		return nil, ae
	case 404:
		return nil, ErrNotFound
	}

	ps := &ProcessStatus{}
	err = json.Unmarshal(buf, ps)
	if err != nil {
		return nil, err
	}

	return ps, nil
}

// GetPDF gets the PDF from the server. In case of an error, the bytes written to w
// might not contain a valid PDF.
func (p *PublishResponse) GetPDF(w io.Writer) error {
	loc := p.endpoint.location + "/pdf/" + p.ID
	req, err := http.NewRequest("GET", loc, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(p.endpoint.password, "")

	resp, err := p.endpoint.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf(resp.Status)
	}
	_, err = io.Copy(w, resp.Body)
	return err
}

// Publish sends data to the server.
func (e *Endpoint) Publish(data *PublishRequest) (PublishResponse, error) {
	var p PublishResponse
	loc := e.location + "/publish?version=" + data.Version
	b, err := json.Marshal(data)
	if err != nil {
		return p, err
	}
	br := bytes.NewReader(b)

	req, err := http.NewRequest("POST", loc, br)
	if err != nil {
		return p, err
	}
	req.SetBasicAuth(e.password, "")
	req.Header.Add("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return p, err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return p, err
	}
	if resp.StatusCode != 201 {
		var ae Error
		err = json.Unmarshal(buf, &ae)
		if err != nil {
			return p, err
		}
		return p, ae
	}

	err = json.Unmarshal(buf, &p)
	if err != nil {
		return p, err
	}
	p.endpoint = e

	return p, nil
}

// NewPublishRequest is the base structure to start a publishing request to the
// endpoint.
func (e *Endpoint) NewPublishRequest() *PublishRequest {
	p := &PublishRequest{
		endpoint: e,
		Version:  "latest",
	}
	return p
}
