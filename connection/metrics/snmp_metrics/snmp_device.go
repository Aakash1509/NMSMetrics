package snmp_metrics

import "Plugin/connection/clients"

// SNMPDeviceMetrics holds the SNMP device metrics
type SNMPDeviceMetrics struct {
	StartedTime    interface{} `json:"started.time"`
	SystemName     interface{} `json:"system.name"`
	SystemLocation interface{} `json:"system.location"`
	SystemDesc     interface{} `json:"system.description"`
}

// GetSNMPDeviceMetrics fetches the SNMP device metrics
func GetSNMPDeviceMetrics(snmpClient *clients.SNMPClient) (*SNMPDeviceMetrics, error) {
	metrics := &SNMPDeviceMetrics{}

	// Define the OIDs to fetch the metrics
	oids := map[string]*interface{}{
		"1.3.6.1.2.1.1.3.0": &metrics.StartedTime,
		"1.3.6.1.2.1.1.1.0": &metrics.SystemName,
		"1.3.6.1.2.1.1.6.0": &metrics.SystemLocation,
		"1.3.6.1.2.1.1.2.0": &metrics.SystemDesc,
	}

	for oid, field := range oids {
		result, err := snmpClient.Client.Get([]string{oid})
		if err != nil {
			*field = "ERROR"
			continue // Continue with the next OID
		}
		for _, variable := range result.Variables {
			*field = variable.Value
		}
	}
	return metrics, nil
}
