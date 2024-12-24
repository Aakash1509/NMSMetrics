package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"strings"
	"time"
)

type SSHClient struct {
	client *ssh.Client
}

func connectSSH(ip string, port int, userName, userPassword string) (*SSHClient, error) {

	config := &ssh.ClientConfig{
		User: userName,

		Auth: []ssh.AuthMethod{

			ssh.Password(userPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		Timeout: 60 * time.Second,
	}

	// Connect to the SSH server
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH server: %v", err)
	}

	return &SSHClient{client: conn}, nil
}

func (s *SSHClient) RunCommand(cmd string) (string, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.Output(cmd)
	return string(output), err
}

func DiscoverSSH(ip string, port int, credentials []Credential) (string, error) {
	for _, profile := range credentials {
		if client, err := connectSSH(ip, port, profile.Name, profile.Password); err == nil {
			defer client.client.Close()

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

func PollSSH(ip, metricGroup, userName, userPassword string) (string, error) {
	sshClient, err := connectSSH(ip, 22, userName, userPassword)
	if err != nil {
		return "", err
	}
	defer sshClient.client.Close()

	// Now, based on the metric group, fetch the appropriate metrics
	var data interface{}

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

	output, _ := json.MarshalIndent(data, "", "  ")
	// Output the formatted JSON
	return string(output), nil
}
