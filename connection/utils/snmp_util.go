package utils

import (
	"Plugin/connection/clients"
	"Plugin/connection/metrics/snmp_metrics"
	"Plugin/connection/models"
	"encoding/json"
	"fmt"
	"github.com/gosnmp/gosnmp"
)

func DiscoverSNMP(input models.Input) (string, error) {

	for _, profile := range input.Credentials {
		snmpClient, err := clients.ConnectSNMP(input.IP, input.Port, profile.Community, profile.Version)
		if err != nil {

			continue //go to next profile
		}
		defer snmpClient.Client.Conn.Close()

		// Hostname OID (sysName)
		oid := []string{"1.3.6.1.2.1.1.5.0"}

		// Fetch the SNMP data
		result, err := snmpClient.Client.Get(oid)
		if err != nil {
			continue
		}

		// Traverse through the SNMP variables to fetch the hostname
		for _, variable := range result.Variables {
			if variable.Type == gosnmp.OctetString {
				// Marshal result into JSON and return it
				result := map[string]interface{}{
					"status":                "Up",
					"hostname":              string(variable.Value.([]byte)),
					"credential.profile.id": profile.ID,
				}
				outputJSON, err := json.Marshal(result)
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
	outputJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %v", err)
	}
	return string(outputJSON), nil
}

func PollSNMP(input models.Input) (string, error) {
	snmpClient, err := clients.ConnectSNMP(input.IP, 161, input.Community, input.Version)
	if err != nil {
		return "", err
	}
	defer snmpClient.Client.Conn.Close()

	var data interface{}
	switch input.MetricGroup {
	case "SNMP.Device":
		data, err = snmp_metrics.GetSNMPDeviceMetrics(snmpClient)
	case "SNMP.Interface":
		data, err = snmp_metrics.GetSNMPInterfaceMetrics(snmpClient)
	default:
		return "", fmt.Errorf("invalid metric group: %s", input.MetricGroup)
	}

	if err != nil {
		return "", err
	}

	output, _ := json.MarshalIndent(data, "", "  ")

	return string(output), nil
}
