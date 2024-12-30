package clients

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
)

type SNMPClient struct {
	Client *gosnmp.GoSNMP
}

func ConnectSNMP(ip string, port int, community, version string) (*SNMPClient, error) {

	// Set up the SNMP Client configuration
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

	return &SNMPClient{Client: params}, nil
}
