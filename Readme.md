# Sample API library for speedata Publisher service // Go

This Go library connects to the speedata api (still in closed beta) and gets PDF from the server.

## Example usage

```go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	api "github.com/speedata/publisher-api"
)

func dothings() error {
	var err error

	ep, err := api.NewEndpoint("username", "https://api.speedata.de")
	if err != nil {
		return err
	}

	p := ep.NewPublishRequest()
	p.AttachFile(filepath.Join("sample", "layout.xml"))
	p.AttachFile(filepath.Join("sample", "data.xml"))
	fmt.Println("-> Now sending data to the server")
	resp, err := ep.Publish(p)
	if err != nil {
		return err
	}

	fmt.Println("-> Getting the status of our publishing run")
	ps, err := resp.Status()
	if err != nil {
		return err
	}
	if ps.Finished != nil {
		fmt.Println("PDF done", ps.Errors, "errors occured")
		for _, e := range ps.Errormessages {
			fmt.Println("*  message", e.Error)
		}
	} else {
		fmt.Println("PDF not finished yet")
	}

	fmt.Println("-> Waiting for the PDF to get written.")
	ps, err = resp.Wait()
	fmt.Println("PDF done", ps.Errors, "errors occured. Finished at", ps.Finished.Format(time.Stamp))
	for _, e := range ps.Errormessages {
		fmt.Println("*  message", e.Error)
	}

	fmt.Println("-> Getting the PDF")
	f, err := os.Create("out.pdf")
	if err != nil {
		return err
	}
	defer f.Close()
	err = resp.GetPDF(f)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := dothings()
	if err != nil {
		if apierror, ok := err.(api.Error); ok {
			fmt.Println("API error", apierror.ErrorType)
			fmt.Println("Instance", apierror.Instance)
			fmt.Println("Title", apierror.Title)
			fmt.Println("Detail", apierror.Detail)
			fmt.Println("Request ID", apierror.RequestID)
		}
		log.Fatal(err)
	}
}
```


Expect changes. This is the very first draft.