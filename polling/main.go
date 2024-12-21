package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Input struct {
	IP          string `json:"ip"`
	MetricGroup string `json:"metric.group.name"`
	Protocol    string `json:"profile.protocol"`
	Name        string `json:"user.name"`
	Password    string `json:"user.password"`
	Community   string `json:"community"`
	Version     string `json:"version"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: <executable> <Input_JSON>")
		os.Exit(1)
	}

	var input Input
	if err := json.Unmarshal([]byte(os.Args[1]), &input); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Connect to SSH
	sshClient, err := ConnectSSH(input.IP, input.Name, input.Password)
	if err != nil {
		log.Fatalf("Error connecting to SSH: %v", err)
	}
	defer sshClient.client.Close()

	var data interface{}

	switch input.MetricGroup {
	case "Linux.Device":
		data, err = GetDeviceMetrics(sshClient)

	case "Linux.CPU":
		data, err = GetCPUMetrics(sshClient)

	case "Linux.Disk":
		data, err = GetDiskMetrics(sshClient)

	case "Linux.Process":
		data, err = GetProcessMetrics(sshClient)
	default:
		log.Fatalf("Invalid metric group: %s", input.MetricGroup)
	}

	if err != nil {
		log.Fatalf("Error fetching metrics: %v", err)
	}

	output, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(output))
}
