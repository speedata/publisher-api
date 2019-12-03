# Sample API library for speedata Publisher service // Go

This Go library connects to the speedata api (still in closed beta) and gets PDF from the server.

## Example usage

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

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
	_, err = toml.DecodeFile("config.toml", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	ep, err := api.NewEndpoint(cfg.Username, cfg.ServerAddress)
	if err != nil {
		return err
	}

	p := ep.NewPublishRequest()
	p.AttachFile(filepath.Join("sample", "layout.xml"))
	p.AttachFile(filepath.Join("sample", "data.xml"))
	fmt.Println("Now sending data to the server")
	resp, err := ep.Publish(p)
	if err != nil {
		return err
	}
	fmt.Println("Getting the PDF")
	b, err := resp.GetPDF()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("out.pdf", b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := dothings()
	if err != nil {
		log.Fatal(err)
	}
}
```


Expect changes. This is the very first draft.