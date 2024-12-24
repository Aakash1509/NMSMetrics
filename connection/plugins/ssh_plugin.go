package plugins

import (
	"Plugin/connection/models"
	"Plugin/connection/utils"
	"fmt"
	"log"
)

func SSHPlugin(input models.Input) (interface{}, error) {

	switch input.EventType {

	case "discover":
		result, err := utils.DiscoverSSH(input)
		if err != nil {
			log.Printf("DiscoverSSH error: %v", err)

			return nil, err
		}
		return result, nil

	case "poll":
		result, err := utils.PollSSH(input)
		if err != nil {
			log.Printf("Error in polling: %v", err)

			return nil, err
		}
		return result, nil

	default:
		return nil, fmt.Errorf("invalid event type: %s", input.EventType) //invalid event type
	}
}
