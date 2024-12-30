package plugins

import (
	"Plugin/metrics/clients"
	"fmt"
	"log"
	"strings"
)

func SSHPlugin(input map[string]interface{}) (interface{}, error) {

	eventType := input["event.type"].(string)

	switch eventType {

	case "discover":
		// Call the utility functions dynamically
		result, err := DiscoverSSH(input)
		if err != nil {
			log.Printf("DiscoverSSH error: %v", err)
			return nil, err
		}
		return result, nil

	case "poll":
		result, err := PollSSH(input)
		if err != nil {
			log.Printf("Error in polling: %v", err)
			return nil, err
		}
		return result, nil

	default:
		return nil, fmt.Errorf("invalid event type: %s", eventType) //invalid event type
	}
}

func DiscoverSSH(input map[string]interface{}) (map[string]interface{}, error) {

	ip := input["ip"].(string)

	port := input["port"].(float64)

	credentials := input["discovery.credential.profiles"].([]interface{})

	for _, credential := range credentials {

		profile := credential.(map[string]interface{})

		name := profile["user_name"].(string)

		password := profile["user_password"].(string)

		client, err := clients.ConnectSSH(ip, int(port), name, password)

		if err == nil {
			defer client.Client.Close()

			output, err := client.RunCommand("hostname")

			if err == nil {

				input["status"] = "Up"
				input["hostname"] = strings.TrimSpace(output)
				input["credential.profile.id"] = profile["profile_id"]
				return input, nil
			}
		}
	}
	input["status"] = "Down"
	input["hostname"] = nil
	input["credential.profile.id"] = nil
	return input, nil
}

func PollSSH(input map[string]interface{}) (map[string]interface{}, error) {

	ip := input["ip"].(string)

	name := input["user.name"].(string)

	password := input["user.password"].(string)

	sshClient, err := clients.ConnectSSH(ip, 22, name, password)

	if err != nil {
		return nil, err
	}

	defer sshClient.Client.Close()

	// Now, based on the metric group, fetch the appropriate metrics
	var data interface{}

	metricGroup := input["metric.group.name"].(string)

	switch metricGroup {

	case "Linux.Device":

		data, err = GetDeviceMetrics(sshClient)

	case "Linux.CPU":

		data, err = GetCPUMetrics(sshClient)

	case "Linux.Disk":

		data, err = GetDiskMetrics(sshClient)

	case "Linux.Process":

		data, err = GetProcessMetrics(sshClient)

	default:

		log.Fatalf("Invalid metric group: %s", metricGroup)
	}

	if err != nil {
		return nil, err
	}

	// Populate the result in the input map
	input["result"] = map[string]interface{}{
		metricGroup: data,
	}
	return input, nil
}
