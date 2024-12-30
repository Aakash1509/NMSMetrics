package plugins

import (
	"Plugin/metrics/clients"
	"fmt"
	"github.com/gosnmp/gosnmp"
	"log"
)

func SNMPPlugin(input map[string]interface{}) (interface{}, error) {

	eventType := input["event.type"].(string)

	switch eventType {

	case "discover":
		result, err := DiscoverSNMP(input)
		if err != nil {
			log.Printf("DiscoverSSH error: %v", err)

			return nil, err
		}
		return result, nil

	case "poll":
		result, err := PollSNMP(input)
		if err != nil {
			log.Printf("Error in polling: %v", err)

			return nil, err
		}
		return result, nil

	default:
		return nil, fmt.Errorf("invalid event type: %s", eventType) //invalid event type
	}
}

func DiscoverSNMP(input map[string]interface{}) (map[string]interface{}, error) {

	ip := input["ip"].(string)

	port := input["port"].(float64)

	credentials := input["discovery.credential.profiles"].([]interface{})

	for _, credential := range credentials {

		profile := credential.(map[string]interface{})

		snmpClient, err := clients.ConnectSNMP(ip, int(port), profile["community"].(string), profile["version"].(string))

		if err == nil {
			defer snmpClient.Client.Conn.Close()

			// Hostname OID (sysName)
			oid := []string{"1.3.6.1.2.1.1.5.0"}

			// Fetch the SNMP data
			result, err := snmpClient.Client.Get(oid)

			if err == nil {
				for _, variable := range result.Variables {
					if variable.Type == gosnmp.OctetString {
						input["status"] = "Up"
						input["hostname"] = string(variable.Value.([]byte))
						input["credential.profile.id"] = profile["profile_id"]
						return input, nil
					}
				}
			}
		}
	}

	input["status"] = "Down"
	input["hostname"] = nil
	input["credential.profile.id"] = nil

	return input, nil
}

func PollSNMP(input map[string]interface{}) (map[string]interface{}, error) {

	ip := input["ip"].(string)

	port := input["port"].(float64)

	community := input["community"].(string)

	version := input["version"].(string)

	snmpClient, err := clients.ConnectSNMP(ip, int(port), community, version)

	if err != nil {
		return nil, err
	}

	defer snmpClient.Client.Conn.Close()

	var data interface{}

	metricGroup := input["metric.group.name"].(string)

	switch metricGroup {
	case "SNMP.Device":
		data, err = GetSNMPDeviceMetrics(snmpClient)
	case "SNMP.Interface":
		data, err = GetSNMPInterfaceMetrics(snmpClient)
	default:

		log.Fatalf("Invalid metric group: %s", metricGroup)
	}

	if err != nil {
		return nil, err
	}

	input["result"] = map[string]interface{}{
		metricGroup: data,
	}
	return input, nil
}
