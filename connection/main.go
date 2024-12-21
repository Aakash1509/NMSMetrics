package main

import (
	"encoding/json"
	"fmt"
	"github.com/gosnmp/gosnmp"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"strconv"
)

type CredentialProfile struct {
	ID        int64  `json:"profile.id"`
	Protocol  string `json:"profile.protocol"`
	Username  string `json:"user.name"`
	Password  string `json:"user.password"`
	Community string `json:"community"`
	Version   string `json:"version"`
}
type DeviceInfo struct {
	IPAddress          string              `json:"discovery.ip"`
	Port               string              `json:"discovery.port"`
	CredentialProfiles []CredentialProfile `json:"discovery.credential.profiles"`
}

func main() {
	log.SetOutput(os.Stderr)

	if len(os.Args) < 4 {
		log.Fatal("Usage: ./program <ip> <port> <credential_profiles_json>")
	}

	ip := os.Args[1]
	port := os.Args[2]
	credentialProfiles := os.Args[3]

	var deviceInfo DeviceInfo

	err := json.Unmarshal([]byte(credentialProfiles), &deviceInfo.CredentialProfiles)
	if err != nil {
		log.Fatal("Failed to parse credential profiles JSON:", err)
	}

	deviceInfo.IPAddress = ip
	deviceInfo.Port = port

	for _, profile := range deviceInfo.CredentialProfiles {
		switch profile.Protocol {
		case "SSH":
			hostname, err := sshConnection(deviceInfo.IPAddress, deviceInfo.Port, profile)
			if err == nil {
				result := map[string]interface{}{
					"status":                "Up",
					"credential.profile.id": profile.ID,
					"hostname":              hostname,
				}
				output, _ := json.Marshal(result)
				fmt.Println(string(output))
				return
			}
		case "SNMP":
			hostname, err := snmpConnection(deviceInfo.IPAddress, deviceInfo.Port, profile)
			if err == nil {
				result := map[string]interface{}{
					"status":                "Up",
					"credential.profile.id": profile.ID,
					"hostname":              hostname,
				}
				output, _ := json.Marshal(result)
				fmt.Println(string(output))
				return
			}
		default:
			log.Println("Unknown protocol:", profile.Protocol)
		}
	}

	// If no valid credential, write to stdout
	result := map[string]interface{}{
		"status":                "Down",
		"credential.profile.id": nil,
		"hostname":              nil,
	}
	output, _ := json.Marshal(result)
	fmt.Println(string(output))
}

// SSH function
func sshConnection(ip, port string, profile CredentialProfile) (string, error) {
	config := &ssh.ClientConfig{
		User: profile.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(profile.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the SSH server
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", ip, port), config)
	if err != nil {
		return "", fmt.Errorf("failed to connect to SSH server: %w", err)
	}
	defer client.Close() // Close the connection

	// Create a new SSH session
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	// Run the command to fetch the hostname
	output, err := session.CombinedOutput("hostname")
	if err != nil {
		return "", fmt.Errorf("failed to run hostname command: %w", err)
	}

	return string(output), nil
}

// SNMP function
func snmpConnection(ip, port string, profile CredentialProfile) (string, error) {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return "", fmt.Errorf("invalid port number: %w", err)
	}

	// Configure SNMP parameters
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

	// OID to get the system hostname
	oids := []string{"1.3.6.1.2.1.1.5.0"}

	// Perform SNMP Get
	result, err := params.Get(oids)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve hostname via SNMP: %w", err)
	}

	// Process the SNMP response
	for _, variable := range result.Variables {
		switch variable.Type {
		case gosnmp.OctetString:
			return string(variable.Value.([]byte)), nil
		default:
			return "", fmt.Errorf("unexpected SNMP type: %v", variable.Type)
		}
	}

	return "", fmt.Errorf("no valid response received from SNMP")
}
