package utils

import (
	"Plugin/connection/clients"
	"Plugin/connection/metrics/ssh_metrics"
	"Plugin/connection/models"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func DiscoverSSH(input models.Input) (string, error) {
	for _, profile := range input.Credentials {
		if client, err := clients.ConnectSSH(input.IP, input.Port, profile.Name, profile.Password); err == nil {
			defer client.Client.Close()

			// Run the hostname command
			if output, err := client.RunCommand("hostname"); err == nil {
				result := map[string]interface{}{
					"status":                "Up",
					"hostname":              strings.TrimSpace(output),
					"credential.profile.id": profile.ID,
				}
				// Marshal result into JSON and return it
				outputJSON, err := json.Marshal(result)
				if err != nil {
					return "", fmt.Errorf("failed to marshal result: %v", err)
				}
				return string(outputJSON), nil
			}
		}
	}

	// If no valid profile is found
	result := map[string]interface{}{
		"status":                "Down",
		"hostname":              nil,
		"credential.profile.id": nil,
	}
	outputJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %v", err)
	}
	return string(outputJSON), nil
}

func PollSSH(input models.Input) (string, error) {
	sshClient, err := clients.ConnectSSH(input.IP, 22, input.UserName, input.UserPassword)
	if err != nil {
		return "", err
	}
	defer sshClient.Client.Close()

	// Now, based on the metric group, fetch the appropriate metrics
	var data interface{}

	switch input.MetricGroup {

	case "Linux.Device":

		data, err = ssh_metrics.GetDeviceMetrics(sshClient)

	case "Linux.CPU":

		data, err = ssh_metrics.GetCPUMetrics(sshClient)

	case "Linux.Disk":

		data, err = ssh_metrics.GetDiskMetrics(sshClient)

	case "Linux.Process":

		data, err = ssh_metrics.GetProcessMetrics(sshClient)

	default:

		log.Fatalf("Invalid metric group: %s", input.MetricGroup)
	}

	output, _ := json.MarshalIndent(data, "", "  ")
	// Output the formatted JSON
	return string(output), nil
}
