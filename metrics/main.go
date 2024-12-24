package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Credential represents a credential profile
type Credential struct {
	ID        int64  `json:"profile_id"`
	Protocol  string `json:"profile_protocol"`
	Name      string `json:"user_name"`
	Password  string `json:"user_password"`
	Community string `json:"community"`
	Version   string `json:"version"`
}

// Input represents the structure of the input JSON
type Input struct {
	EventType       string       `json:"event.type"`  // Either "discover" or "poll"
	DeviceType      string       `json:"device.type"` // Either "Linux" or "SNMP"
	IP              string       `json:"ip"`
	Port            int          `json:"port"`                          // Used for "discover"
	Credentials     []Credential `json:"discovery.credential.profiles"` // Used for "discover"
	MetricGroup     string       `json:"metric.group.name"`             // Used for "poll"
	ProfileProtocol string       `json:"profile.protocol"`              // Used for "poll"
	UserName        string       `json:"user.name"`                     // Used for "poll"
	UserPassword    string       `json:"user.password"`                 // Used for "poll"
	Community       string       `json:"community"`                     // Used for "poll"
	Version         string       `json:"version"`                       // Used for "poll"
}

func main() {
	if len(os.Args) != 2 {

		fmt.Println("Usage: <executable> <Input_JSON>")

		os.Exit(1)
	}

	var input Input

	err := json.Unmarshal([]byte(os.Args[1]), &input)

	if err != nil {

		log.Fatalf("Error parsing JSON: %v", err)
	}

	var data interface{}

	switch input.DeviceType {
	case "Linux":
		data, err = SSHPlugin(input)
	case "SNMP":
		data, err = SNMPPlugin(input)
	default:
		log.Fatalf("Unsupported device type: %s", input.DeviceType)
	}

	if err != nil {
		log.Fatalf("Error executing plugin: %v", err)
	}

	// Print the result (JSON formatted string)
	fmt.Println(data)
}
