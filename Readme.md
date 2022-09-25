[![Go Reference](https://pkg.go.dev/badge/github.com/speedata/publisher-api.svg)](https://pkg.go.dev/github.com/speedata/publisher-api)


# Sample API library for speedata Publisher PDF generation service // Go

This Go library connects to the speedata api and gets PDF from the server. See https://doc.speedata.de/publisher/en/saasapi/ for a description of the API.


## Example usage

```go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	api "github.com/speedata/publisher-api"
)

type config struct {
	ServerAddress string
	Username      string
}

func dothings() error {
	var err error
	cfg := config{}
	_, err = toml.DecodeFile("clientconfig.toml", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	ep, err := api.NewEndpoint(cfg.Username, cfg.ServerAddress)
	if err != nil {
		return err
	}
	p := ep.NewPublishRequest()

	// If you need a different specific version, you can check the availability
	// here:
	// versions, err := ep.AvailableVersions()
	// if err != nil {
	//    return err
	// }
	// and set the required version
	// p.Version = versions[0]

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
	if err != nil {
		return err
	}

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

