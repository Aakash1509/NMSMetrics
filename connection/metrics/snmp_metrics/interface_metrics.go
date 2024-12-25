package snmp_metrics

import (
	"Plugin/connection/clients"
	"strconv"
)

// SNMPInterfaceMetrics holds the SNMP interface metrics
type SNMPInterfaceMetrics struct {
	InterfaceSentDiscardPackets     map[int]interface{} `json:"interface.sent.discard.packets"`
	InterfaceInPackets              map[int]interface{} `json:"interface.in.packets"`
	InterfacePackets                map[int]interface{} `json:"interface.packets"`
	InterfaceErrorPackets           map[int]interface{} `json:"interface.error.packets"`
	InterfaceSentErrorPackets       map[int]interface{} `json:"interface.sent.error.packets"`
	InterfaceReceivedDiscardPackets map[int]interface{} `json:"interface.received.discard.packets"`
	InterfaceReceivedOctets         map[int]interface{} `json:"interface.received.octets"`
	InterfaceOutPackets             map[int]interface{} `json:"interface.out.packets"`
	InterfaceOperationalStatus      map[int]interface{} `json:"interface.operational.status"`
	InterfaceAdminStatus            map[int]interface{} `json:"interface.admin.status"`
	InterfaceReceivedErrorPackets   map[int]interface{} `json:"interface.received.error.packets"`
	InterfaceDiscardPacket          map[int]interface{} `json:"interface.discard.packet"`
}

// GetSNMPInterfaceMetrics fetches the SNMP interface metrics
func GetSNMPInterfaceMetrics(snmpClient *clients.SNMPClient) (*SNMPInterfaceMetrics, error) {
	metrics := &SNMPInterfaceMetrics{
		InterfaceSentDiscardPackets:     make(map[int]interface{}),
		InterfaceInPackets:              make(map[int]interface{}),
		InterfacePackets:                make(map[int]interface{}),
		InterfaceErrorPackets:           make(map[int]interface{}),
		InterfaceSentErrorPackets:       make(map[int]interface{}),
		InterfaceReceivedDiscardPackets: make(map[int]interface{}),
		InterfaceReceivedOctets:         make(map[int]interface{}),
		InterfaceOutPackets:             make(map[int]interface{}),
		InterfaceOperationalStatus:      make(map[int]interface{}),
		InterfaceAdminStatus:            make(map[int]interface{}),
		InterfaceReceivedErrorPackets:   make(map[int]interface{}),
		InterfaceDiscardPacket:          make(map[int]interface{}),
	}

	// Define the base OIDs for the metrics
	baseOids := map[string]string{
		"sent.discard.packets":     "1.3.6.1.2.1.2.2.1.10", // ifOutDiscard
		"in.packets":               "1.3.6.1.2.1.2.2.1.3",  // ifInUcastPkts
		"packets":                  "1.3.6.1.2.1.2.2.1.7",  // ifOperStatus
		"error.packets":            "1.3.6.1.2.1.2.2.1.13", // ifInErrors
		"sent.error.packets":       "1.3.6.1.2.1.2.2.1.14", // ifOutErrors
		"received.discard.packets": "1.3.6.1.2.1.2.2.1.15", // ifInDiscards
		"received.octets":          "1.3.6.1.2.1.2.2.1.16", // ifInOctets
		"out.packets":              "1.3.6.1.2.1.2.2.1.17", // ifOutOctets
		"operational.status":       "1.3.6.1.2.1.2.2.1.19", // ifAdminStatus
		"admin.status":             "1.3.6.1.2.1.2.2.1.8",  // ifAdminStatus
		"received.error.packets":   "1.3.6.1.2.1.2.2.1.18", // ifInErrors
		"discard.packet":           "1.3.6.1.2.1.2.2.1.25", // ifDiscards
	}

	// Iterate over the base OIDs and fetch the data for each interface index
	for metricName, baseOid := range baseOids {
		// Fetch the interface count first
		interfaceCountOid := "1.3.6.1.2.1.2.1.0" // ifNumber OID to get the number of interfaces
		result, err := snmpClient.Client.Get([]string{interfaceCountOid})
		if err != nil {
			return nil, err
		}

		// Get the number of interfaces
		var interfaceCount int
		for _, variable := range result.Variables {
			interfaceCount = int(variable.Value.(int))
		}

		// Fetch the data for each interface
		for i := 1; i <= interfaceCount; i++ {
			oid := baseOid + "." + strconv.Itoa(i)
			result, err := snmpClient.Client.Get([]string{oid})
			if err != nil {
				return nil, err
			}

			// Assign the value to the correct metric map
			switch metricName {
			case "sent.discard.packets":
				metrics.InterfaceSentDiscardPackets[i] = result.Variables[0].Value
			case "in.packets":
				metrics.InterfaceInPackets[i] = result.Variables[0].Value
			case "packets":
				metrics.InterfacePackets[i] = result.Variables[0].Value
			case "error.packets":
				metrics.InterfaceErrorPackets[i] = result.Variables[0].Value
			case "sent.error.packets":
				metrics.InterfaceSentErrorPackets[i] = result.Variables[0].Value
			case "received.discard.packets":
				metrics.InterfaceReceivedDiscardPackets[i] = result.Variables[0].Value
			case "received.octets":
				metrics.InterfaceReceivedOctets[i] = result.Variables[0].Value
			case "out.packets":
				metrics.InterfaceOutPackets[i] = result.Variables[0].Value
			case "operational.status":
				metrics.InterfaceOperationalStatus[i] = result.Variables[0].Value
			case "admin.status":
				metrics.InterfaceAdminStatus[i] = result.Variables[0].Value
			case "received.error.packets":
				metrics.InterfaceReceivedErrorPackets[i] = result.Variables[0].Value
			case "discard.packet":
				metrics.InterfaceDiscardPacket[i] = result.Variables[0].Value
			}
		}
	}

	return metrics, nil
}
