package plugins

import (
	"Plugin/metrics/clients"
	"fmt"
	"strconv"
)

var baseOids = map[string]string{
	"interface.sent.discard.packets":     "1.3.6.1.2.1.2.2.1.10",
	"interface.in.packets":               "1.3.6.1.2.1.2.2.1.3",
	"interface.packets":                  "1.3.6.1.2.1.2.2.1.7",
	"interface.error.packets":            "1.3.6.1.2.1.2.2.1.13",
	"interface.sent.error.packets":       "1.3.6.1.2.1.2.2.1.14",
	"interface.received.discard.packets": "1.3.6.1.2.1.2.2.1.15",
	"interface.received.octets":          "1.3.6.1.2.1.2.2.1.16",
	"interface.out.packets":              "1.3.6.1.2.1.2.2.1.17",
	"interface.operational.status":       "1.3.6.1.2.1.2.2.1.19",
	"interface.admin.status":             "1.3.6.1.2.1.2.2.1.8",
	"interface.received.error.packets":   "1.3.6.1.2.1.2.2.1.18",
	"interface.discard.packet":           "1.3.6.1.2.1.2.2.1.25",
}

func GetSNMPInterfaceMetrics(snmpClient *clients.SNMPClient) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})

	interfaceCountOid := "1.3.6.1.2.1.2.1.0" // OID to get the number of interfaces

	result, err := snmpClient.Client.Get([]string{interfaceCountOid})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch interface count: %w", err)
	}

	// Number of interfaces
	var interfaceCount int

	for _, variable := range result.Variables {
		switch v := variable.Value.(type) {
		case int:
			interfaceCount = v
		case uint32:
			interfaceCount = int(v)
		default:
			return nil, fmt.Errorf("unexpected type for interface count: %T", v)
		}
	}

	// Fetch metrics for each interface
	for metricName, baseOid := range baseOids {
		metricData := make(map[string]interface{})
		for i := 1; i <= interfaceCount; i++ {
			oid := baseOid + "." + strconv.Itoa(i)
			result, err := snmpClient.Client.Get([]string{oid})
			if err != nil {
				continue
			}

			// Assign the value for the specific interface index
			if len(result.Variables) > 0 {
				variable := result.Variables[0]
				if variable.Value != nil {
					metricData[strconv.Itoa(i)] = variable.Value
				}
			}

		}

		if len(metricData) > 0 {
			metrics[metricName] = metricData
		}
	}
	return metrics, nil
}
