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

const (
	API_PREFIX string = "/v0"
)

func basicAuth(username string) string {
	auth := username + ":"
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// PublishRequest is
type PublishRequest struct {
	endpoint *Endpoint
	Files    []PublishFile `json:"files"`
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
	Filename string `json:"filename"`
	Contents []byte `json:"contents"`
}

// PublishResponse holds the id to the publishing process
type PublishResponse struct {
	endpoint *Endpoint
	Id       string
}

type Errormessage struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type ProcessStatus struct {
	Finished      *time.Time
	Errors        int
	Errormessages []Errormessage
}

func (p *PublishResponse) Status() (*ProcessStatus, error) {
	loc := p.endpoint.location + "/status/" + p.Id
	req, err := http.NewRequest("GET", loc, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+p.endpoint.password)

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

// Wait for the publishing process to finish. Return an error if something is wrong with the request.
// If there is an error during the publishing run but the request itself is without errors, the
// error is nil, but the returned publishing status has the numbers of errors et.
func (p *PublishResponse) Wait() (*ProcessStatus, error) {
	loc := p.endpoint.location + "/wait/" + p.Id
	req, err := http.NewRequest("GET", loc, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+p.endpoint.password)

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

// GetPDF gets the PDF from the server. In case of an error, the byte slice might not be meaningful.
// Otherwise it holds the PDF file.
func (p *PublishResponse) GetPDF(w io.Writer) error {
	loc := p.endpoint.location + "/pdf/" + p.Id
	req, err := http.NewRequest("GET", loc, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Basic "+p.endpoint.password)

	resp, err := p.endpoint.client.Do(req)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, resp.Body)
	return err
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
		fmt.Println("new request failed")
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

	if resp.StatusCode != 201 {
		fmt.Println(string(buf))
		var ae APIError
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

// NewPublishRequest is the base structure to start a publishing request to the endpoint.
func (e *Endpoint) NewPublishRequest() *PublishRequest {
	p := &PublishRequest{
		endpoint: e,
	}
	return p
}
