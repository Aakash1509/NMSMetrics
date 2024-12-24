package ssh_metrics

import (
	"Plugin/connection/clients"
	"Plugin/connection/utils/parse"
	"strings"
)

type DiskMetrics struct {
	Volume        interface{} `json:"system.disk.volume"`
	MountPath     interface{} `json:"system.disk.volume.mount.path"`
	UsedPercent   interface{} `json:"system.disk.volume.used.percent"`
	CapacityBytes interface{} `json:"system.disk.volume.capacity.bytes"`
	FreePercent   interface{} `json:"system.disk.volume.free.percent"`
	UsedBytes     interface{} `json:"system.disk.volume.used.bytes"`
	FreeBytes     interface{} `json:"system.disk.volume.free.bytes"`
}

func GetDiskMetrics(sshClient *clients.SSHClient) ([]DiskMetrics, error) {
	// Command to fetch disk information in a single run
	cmd := `df -B1 --output=source,target,pcent,size,used,avail | tail -n +2`
	result, err := sshClient.RunCommand(cmd)
	if err != nil {
		return nil, err
	}

	// Split the output into lines
	lines := strings.Split(strings.TrimSpace(result), "\n")
	var metricsList []DiskMetrics

	// Process each line of the command output
	for _, line := range lines {
		fields := strings.Fields(line)

		metrics := DiskMetrics{
			Volume:        parse.ParseResult(fields[0]),
			MountPath:     parse.ParseResult(fields[1]),
			UsedPercent:   parse.ParseResult(strings.TrimSuffix(fields[2], "%")), // Remove trailing '%'
			CapacityBytes: parse.ParseResult(fields[3]),
			UsedBytes:     parse.ParseResult(fields[4]),
			FreeBytes:     parse.ParseResult(fields[5]),
			FreePercent:   100 - parse.ParseResult(strings.TrimSuffix(fields[2], "%")).(float64), // I need to type assert as return type is interface{}
		}

		metricsList = append(metricsList, metrics)
	}

	return metricsList, nil
}
