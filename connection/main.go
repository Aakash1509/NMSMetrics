package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gosnmp/gosnmp"
	"golang.org/x/crypto/ssh"
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

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: <executable> <IP> <Port> <CredentialProfiles_JSON>")
		os.Exit(1)
	}

	ip := os.Args[1]

	port := os.Args[2]

	credentialJSON := os.Args[3]

	var profiles []Credential

	err := json.Unmarshal([]byte(credentialJSON), &profiles)
	if err != nil {
		log.Fatalf("Invalid credential profiles JSON: %v", err)
	}

	//Creating a buffered channel with length equal to number of credential profiles
	results := make(chan map[string]interface{}, len(profiles))

	done := make(chan struct{})

	var wg sync.WaitGroup

	for _, profile := range profiles {
		wg.Add(1)

		go func(profile Credential) {

			defer wg.Done()

			var hostname string

			var err error

			if profile.Protocol == "SNMP" {
				hostname, err = snmpConnection(ip, port, profile)
			} else if profile.Protocol == "SSH" {
				hostname, err = sshConnection(ip, port, profile)
			}

			if err == nil {
				select {
				case results <- map[string]interface{}{
					"status":                "Up",
					"credential.profile.id": profile.ID,
					"hostname":              hostname,
				}:
				case <-done:
					return
				}
			}
		}(profile)
	}

	go func() {
		wg.Wait()

		close(results)
	}()

	// Collect results
	for result := range results {
		printResult(
			result["status"].(string),
			result["credential.profile.id"],
			result["hostname"],
		)
		close(done)

		return // Exit on the first successful connection
	}

	// If no match is found
	printResult("Down", nil, nil)
}

func snmpConnection(ip, port string, profile Credential) (string, error) {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return "", fmt.Errorf("invalid port number: %w", err)
	}

	params := &gosnmp.GoSNMP{
		Target:    ip,
		Port:      uint16(portInt),
		Community: profile.Community,
		Version:   gosnmp.Version2c,
		Timeout:   gosnmp.Default.Timeout,
	}

	err = params.Connect()
	if err != nil {
		return "", fmt.Errorf("failed to connect to SNMP server: %w", err)
	}
	defer params.Conn.Close()

	// Get the hostname OID
	oid := []string{"1.3.6.1.2.1.1.5.0"}
	result, err := params.Get(oid)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve SNMP hostname: %w", err)
	}

	for _, variable := range result.Variables {
		if variable.Type == gosnmp.OctetString {
			return string(variable.Value.([]byte)), nil
		}
	}

	return "", fmt.Errorf("SNMP response did not contain valid hostname")
}

func sshConnection(ip, port string, profile Credential) (string, error) {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return "", fmt.Errorf("invalid port number: %w", err)
	}

	config := &ssh.ClientConfig{
		User: profile.Name,
		Auth: []ssh.AuthMethod{
			ssh.Password(profile.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(60) * time.Second,
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, portInt), config)
	if err != nil {
		return "", fmt.Errorf("failed to connect to SSH server: %w", err)
	}
	defer conn.Close()

	// command to fetch the hostname
	session, err := conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	hostname, err := session.Output("hostname")
	if err != nil {
		return "", fmt.Errorf("failed to execute hostname command: %w", err)
	}

	return string(hostname), nil
}

func printResult(status string, profileID interface{}, hostname interface{}) {
	result := map[string]interface{}{
		"status":                status,
		"credential.profile.id": profileID,
		"hostname":              hostname,
	}

	output, _ := json.Marshal(result)
	fmt.Println(string(output))
}
