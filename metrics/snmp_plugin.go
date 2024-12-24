package main

import (
	"fmt"
	"log"
)

func SNMPPlugin(input Input) (interface{}, error) {

	switch input.EventType {

	case "discover":
		result, err := DiscoverSNMP(input.IP, input.Port, input.Credentials)
		if err != nil {
			log.Printf("DiscoverSSH error: %v", err)

			return nil, err
		}
		return result, nil

	case "poll":
		result, err := PollSNMP(input.IP, input.MetricGroup, input.Community, input.Version)
		if err != nil {
			log.Printf("Error in polling: %v", err)

			return nil, err
		}
		return result, nil

	default:
		return nil, fmt.Errorf("invalid event type: %s", input.EventType) //invalid event type
	}
}
