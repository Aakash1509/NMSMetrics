package main

import (
	"encoding/json"
	"fmt"
	"github.com/gosnmp/gosnmp"
)

type SNMPClient struct {
	client *gosnmp.GoSNMP
}

func connectSNMP(ip string, port int, community, version string) (*SNMPClient, error) {

	// Set up the SNMP client configuration
	params := &gosnmp.GoSNMP{

		Target: ip,

		Port: uint16(port),

		Community: community,

		Timeout: gosnmp.Default.Timeout,
	}

	switch version {

	case "v1":
		params.Version = gosnmp.Version1

	case "v2c":
		params.Version = gosnmp.Version2c

	case "v3":
		params.Version = gosnmp.Version3

	default:
		return nil, fmt.Errorf("unsupported SNMP version: %s", version)
	}

	// Connect to the SNMP server
	err := params.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SNMP server: %v", err)
	}

	return &SNMPClient{client: params}, nil
}

func DiscoverSNMP(ip string, port int, credentials []Credential) (string, error) {

	for _, profile := range credentials {
		snmpClient, err := connectSNMP(ip, port, profile.Community, profile.Version)
		if err != nil {

			continue //go to next profile
		}
		defer snmpClient.client.Conn.Close()

		// Hostname OID (sysName)
		oid := []string{"1.3.6.1.2.1.1.5.0"}

		// Fetch the SNMP data
		result, err := snmpClient.client.Get(oid)
		if err != nil {
			continue
		}

		// Traverse through the SNMP variables to fetch the hostname
		for _, variable := range result.Variables {
			if variable.Type == gosnmp.OctetString {
				// Marshal result into JSON and return it
				resultMap := map[string]interface{}{
					"status":                "Up",
					"hostname":              string(variable.Value.([]byte)),
					"credential.profile.id": profile.ID,
				}
				outputJSON, err := json.MarshalIndent(resultMap, "", "  ")
				if err != nil {
					return "", fmt.Errorf("failed to marshal result: %v", err)
				}
				return string(outputJSON), nil
			}
		}
	}

	// If no valid hostname is found after iterating through all credentials
	result := map[string]interface{}{
		"status":                "Down",
		"hostname":              nil,
		"credential.profile.id": nil,
	}
	outputJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %v", err)
	}
	return string(outputJSON), fmt.Errorf("failed to connect to SNMP server with provided credentials")
}

func PollSNMP(ip, metricGroup, community, version string) (string, error) {
	snmpClient, err := connectSNMP(ip, 161, community, version)
	if err != nil {
		return "", err
	}
	defer snmpClient.client.Conn.Close()

	var data interface{}
	switch metricGroup {
	case "SNMP.Device":
		data, err = GetSNMPDeviceMetrics(snmpClient)
	case "SNMP.Interface":
		data, err = GetSNMPInterfaceMetrics(snmpClient)
	default:
		return "", fmt.Errorf("invalid metric group: %s", metricGroup)
	}

	if err != nil {
		return "", err
	}

	output, _ := json.MarshalIndent(data, "", "  ")

	return string(output), nil
}
