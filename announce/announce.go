package announce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func requestServiceAnnouncement(overseerRoot string, service []byte) {
	go func() {
		available := true

		checkchan := make(chan bool)
		errorchan := make(chan error)
		defer func() {
			close(checkchan)
			close(errorchan)
		}()

		for available {

			go func() {
				resp, err := http.Post(overseerRoot, "application/json", bytes.NewBuffer(service))
				if err != nil {
					errorchan <- err
					return
				}
				if resp.StatusCode > 299 {
					errorchan <- fmt.Errorf("Request Error: %s", resp.Status)
					return
				}

				checkchan <- true
			}() // end request routine

			select {
			case <-checkchan:
				// JUST GO ON
			case <-errorchan:
			case <-time.After(time.Second * 3):
			}

			time.Sleep(10 * time.Minute) // update registration once a minute

		} // end for loop

	}() // end loop routine
}

// NewService provides a method to attach a new Service to the overseer stack
func NewService(serviceName string) {
	serviceRoot := os.Getenv("ROOT_URL")
	aliveResource := os.Getenv("ALIVE_URL")
	overseerRoot := os.Getenv("OVERSEER_ROOT")

	if serviceRoot == "" || aliveResource == "" || overseerRoot == "" {
		fmt.Println("ROOT_URL and/or ALIVE_URL and/or OVERSEER_ROOT not set")
		return
	}

	service, _ := json.Marshal(struct {
		Name  string `json:"name"`
		Root  string `json:"root"`
		Alive string `json:"alive"`
	}{
		serviceName,
		serviceRoot,
		aliveResource,
	})

	fmt.Printf("Announcing new Service to Overseer: %s", service)
	requestServiceAnnouncement(overseerRoot, service)
}
