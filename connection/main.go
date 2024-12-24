package main

import (
	"Plugin/connection/models"
	"Plugin/connection/plugins"
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

	var input models.Input

	err := json.Unmarshal([]byte(os.Args[1]), &input)

	if err != nil {

		log.Fatalf("Error parsing JSON: %v", err)
	}

	var data interface{}

	switch input.DeviceType {
	case "Linux":
		data, err = plugins.SSHPlugin(input)
	case "SNMP":
		data, err = plugins.SNMPPlugin(input)
	default:
		log.Fatalf("Unsupported device type: %s", input.DeviceType)
	}

	if err != nil {
		log.Fatalf("Error executing plugin: %v", err)
	}

	// Print the result (JSON formatted string)
	fmt.Println(data)
}
