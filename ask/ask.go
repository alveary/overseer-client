package ask

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Service response from overseer http request
type Service struct {
	Name    string `json:name`
	Address string `json:address`
}

func requestOverseer(name string, errorchan chan error, responsechan chan *http.Response) {
	overseerHost := os.Getenv("OVERSEER_HOST")
	overseerPort := os.Getenv("OVERSEER_PORT")

	resp, err := http.Get(overseerHost + ":" + overseerPort + "/" + name)

	if err != nil {
		errorchan <- err
		return
	}
	if resp.StatusCode > 299 {
		errorchan <- fmt.Errorf("Request Error: %s", resp.Status)
		return
	}

	responsechan <- resp
}

// ForService lets you retrieve the service URI information
func ForService(serviceName string) (retrieved Service, err error) {
	responsechan := make(chan *http.Response)
	errorchan := make(chan error)

	defer func() {
		close(responsechan)
		close(errorchan)
	}()

	go requestOverseer(serviceName, errorchan, responsechan)

	select {
	case resp := <-responsechan:
		defer resp.Body.Close()

		dec := json.NewDecoder(resp.Body)
		dec.Decode(&retrieved)

		return retrieved, nil
	case err = <-errorchan:
		return retrieved, err
	case <-time.After(time.Second * 3):
		return retrieved, fmt.Errorf("Timeout of service registry call")
	}
}
