package plugins

import (
	"Plugin/metrics/clients"
	"fmt"
	"strconv"
)

var snmpOIDs = map[string]string{
	"started.time":       "1.3.6.1.2.1.1.3.0", // Uptime
	"system.name":        "1.3.6.1.2.1.1.5.0", // System name
	"system.location":    "1.3.6.1.2.1.1.6.0", // System location
	"system.description": "1.3.6.1.2.1.1.1.0", // System description
}

func GetSNMPDeviceMetrics(snmpClient *clients.SNMPClient) (map[string]interface{}, error) {

	results := make(map[string]interface{})

	for key, oid := range snmpOIDs {
		response, err := snmpClient.Client.Get([]string{oid})
		if err != nil {
			results[key] = "ERROR"
			continue // Skip to the next OID
		}

		// Process the SNMP response
		for _, variable := range response.Variables {
			switch value := variable.Value.(type) {
			case string:
				results[key] = value
			case uint32:
				results[key] = strconv.Itoa(int(value))
			case []uint8:
				results[key] = string(value)
			default:
				results[key] = fmt.Sprintf("Unsupported type: %T", value) // Handle unexpected types
			}
		}
	}

	return results, nil
}
