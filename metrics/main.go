package main

import (
	"Plugin/metrics/plugins"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: <executable> <Input_JSON>")
		os.Exit(1)
	}

	var input map[string]interface{}

	err := json.Unmarshal([]byte(os.Args[1]), &input)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Extracting required fields
	deviceType, ok := input["device_type"].(string)
	if !ok {
		log.Fatalf("Missing or invalid 'device_type' in input")
	}

	var data interface{}

	switch deviceType {
	case "Linux":
		data, err = plugins.SSHPlugin(input)
	case "SNMP":
		data, err = plugins.SNMPPlugin(input)
	default:
		log.Fatalf("Unsupported device type: %s", deviceType)
	}

	if err != nil {
		log.Fatalf("Error executing plugin: %v", err)
	}

	// Print the result (JSON formatted string)
	result, err := json.Marshal(data)

	if err != nil {
		log.Fatalf("Error marshalling result: %v", err)
	}
	fmt.Println(string(result))
}
