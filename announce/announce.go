package announce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func postService(address string, service []byte, donechan chan bool, errorchan chan error, checkchan chan bool) {
	resp, err := http.Post(address, "application/json", bytes.NewBuffer(service))
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
}

func requestServiceAnnouncement(overseerAddress string, service []byte) {
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

			go postService(overseerAddress, service, donechan, errorchan, checkchan) // end request routine

			select {
			case <-checkchan:
				// JUST GO ON
			case err := <-errorchan:
				log.Println(err)
				time.Sleep(2 * time.Second)
				requestServiceAnnouncement(overseerAddress, service)
				available = false
				return
			case <-time.After(time.Second * 3):
				donechan <- true
				log.Printf("Registering new Service failed with timeout")
				time.Sleep(2 * time.Second)
				requestServiceAnnouncement(overseerAddress, service)
			}

			time.Sleep(10 * time.Minute) // renewal registration once in a time, in case the overseer service failed
		} // end for loop

	}() // end loop routine
}

// Service provides a method to attach a new Service to the overseer stack
func Service(serviceName string) {
	serviceRoot := os.Getenv("ROOT_URL")
	aliveResource := os.Getenv("ALIVE_URL")

	overseerHost := os.Getenv("OVERSEER_HOST")
	overseerPort := os.Getenv("OVERSEER_PORT")
	overseerAddress := overseerHost + ":" + overseerPort

	if serviceRoot == "" || aliveResource == "" || overseerHost == "" || overseerPort == "" {
		log.Println("ROOT_URL and/or ALIVE_URL and/or OVERSEER_ROOT not set")
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
	requestServiceAnnouncement(overseerAddress, service)
}
