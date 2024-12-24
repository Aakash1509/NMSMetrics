package main

// SNMPInterfaceMetrics holds the SNMP interface metrics
type SNMPInterfaceMetrics struct {
	SentDiscardPackets     interface{} `json:"interface.sent.discard.packets"`
	InPackets              interface{} `json:"interface.in.packets"`
	Packets                interface{} `json:"interface.packets"`
	ErrorPackets           interface{} `json:"interface.error.packets"`
	SentErrorPackets       interface{} `json:"interface.sent.error.packets"`
	ReceivedDiscardPackets interface{} `json:"interface.received.discard.packets"`
	ReceivedOctets         interface{} `json:"interface.received.octets"`
	OutPackets             interface{} `json:"interface.out.packets"`
	OperationalStatus      interface{} `json:"interface.operational.status"`
	AdminStatus            interface{} `json:"interface.admin.status"`
	ReceivedErrorPackets   interface{} `json:"interface.received.error.packets"`
	DiscardPacket          interface{} `json:"interface.discard.packet"`
}

// GetSNMPInterfaceMetrics fetches the SNMP interface metrics
func GetSNMPInterfaceMetrics(snmpClient *SNMPClient) (*SNMPInterfaceMetrics, error) {
	metrics := &SNMPInterfaceMetrics{}

	// Define the OIDs to fetch the metrics
	oids := map[string]*interface{}{
		"1.3.6.1.2.1.2.2.1.13": &metrics.SentDiscardPackets,
		"1.3.6.1.2.1.2.2.1.11": &metrics.InPackets,
		"1.3.6.1.2.1.2.2.1.10": &metrics.Packets,
		"1.3.6.1.2.1.2.2.1.19": &metrics.ErrorPackets,
		"1.3.6.1.2.1.2.2.1.21": &metrics.SentErrorPackets,
		"1.3.6.1.2.1.2.2.1.23": &metrics.ReceivedDiscardPackets,
		"1.3.6.1.2.1.2.2.1.16": &metrics.ReceivedOctets,
		"1.3.6.1.2.1.2.2.1.17": &metrics.OutPackets,
		"1.3.6.1.2.1.2.2.1.8":  &metrics.OperationalStatus,
		"1.3.6.1.2.1.2.2.1.7":  &metrics.AdminStatus,
		"1.3.6.1.2.1.2.2.1.20": &metrics.ReceivedErrorPackets,
		"1.3.6.1.2.1.2.2.1.12": &metrics.DiscardPacket,
	}

	for oid, field := range oids {
		result, err := snmpClient.client.Get([]string{oid})
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
