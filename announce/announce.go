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

		donechan := make(chan bool)
		checkchan := make(chan bool)
		errorchan := make(chan error)

		defer func() {
			close(donechan)
			close(checkchan)
			close(errorchan)
		}()

		for available {

			go func() {
				resp, err := http.Post("http://"+overseerRoot, "application/json", bytes.NewBuffer(service))
				select {
				case <-donechan:
					return
				default:
					if err != nil {
						errorchan <- err
						return
					}
					if resp.StatusCode > 299 {
						errorchan <- fmt.Errorf("Request Error: %s", resp.Status)
						return
					}

					checkchan <- true
				}
			}() // end request routine

			select {
			case <-checkchan:
				// JUST GO ON
			case err := <-errorchan:
				fmt.Println(err)
				time.Sleep(2 * time.Second)
				requestServiceAnnouncement(overseerRoot, service)
				available = false
				return
			case <-time.After(time.Second * 3):
				donechan <- true
				fmt.Errorf("Registering new Service failed with timeout")
				time.Sleep(2 * time.Second)
				requestServiceAnnouncement(overseerRoot, service)
			}

			time.Sleep(10 * time.Minute) // renewal registration once in a time, in case the overseer service failed
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
