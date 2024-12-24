package main

import (
	"fmt"
	"log"
)

func SSHPlugin(input Input) (interface{}, error) {

	switch input.EventType {

	case "discover":
		result, err := DiscoverSSH(input.IP, input.Port, input.Credentials)
		if err != nil {
			log.Printf("DiscoverSSH error: %v", err)

			return nil, err
		}
		return result, nil

	case "poll":
		result, err := PollSSH(input.IP, input.MetricGroup, input.UserName, input.UserPassword)
		if err != nil {
			log.Printf("Error in polling: %v", err)

			return nil, err
		}
		return result, nil

	default:
		return nil, fmt.Errorf("invalid event type: %s", input.EventType) //invalid event type
	}
}
