package messaging

import (
	"errors"
	"log"
	"time"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/utils"
)

func waitReply(id string, rchan chan response.RabbitMQResponse) (response.RabbitMQResponse, error) {
	for {
		select {
		case docReply := <-rchan:
			// responses received
			log.Printf("INFO: received reply: %v uid: %s", docReply, id)

			delete(utils.Rchans, id)
			return docReply, nil
		case <-time.After(100 * time.Second):
			// timeout
			log.Printf("ERROR: request timeout uid: %s", id)

			// remove channel from rchans
			delete(utils.Rchans, id)
			return response.RabbitMQResponse{}, errors.New("request timeout")
		}
	}
}
