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

type serviceJSON struct {
	Name string `json:"name"`
	Root string `json:"address"`
}

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
				// TODO: add retry counter to break this circle
				requestServiceAnnouncement(overseerAddress, service)
			}

			time.Sleep(10 * time.Minute) // renewal registration once in a time, in case the overseer service failed
		} // end for loop

	}() // end loop routine
}

func overseerAddress() string {
	overseerHost := os.Getenv("OVERSEER_HOST")
	overseerPort := os.Getenv("OVERSEER_PORT")

	log.Printf("Announce to Overseer: %s", overseerHost)

	if overseerHost == "" || overseerPort == "" {
		log.Fatal("OVERSEER_HOST or OVERSEER_PORT not set")
	}

	return overseerHost + ":" + overseerPort
}

func serviceAddresses() (string, string) {
	serviceRoot := os.Getenv("ROOT_URL")
	aliveResource := os.Getenv("ALIVE_URL")

	if serviceRoot == "" || aliveResource == "" {
		log.Println("ROOT_URL and/or ALIVE_URL")
	}

	return serviceRoot, aliveResource
}

// Service provides a method to attach a new Service to the overseer stack
func Service(serviceName string) {
	serviceRoot, _ := serviceAddresses()
	overseerAddress := overseerAddress()
	service, _ := json.Marshal(serviceJSON{
		serviceName,
		serviceRoot,
	})

	log.Printf("Announcing new Service to Overseer: %s", service)
	requestServiceAnnouncement(overseerAddress, service)
}
