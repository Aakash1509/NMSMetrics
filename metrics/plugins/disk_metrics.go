package plugins

import (
	"Plugin/metrics/clients"
	"Plugin/metrics/utils"
	"strings"
)

var diskCommand = `df -B1 --output=source,target,pcent,size,used,avail | tail -n +2`

func GetDiskMetrics(sshClient *clients.SSHClient) ([]map[string]interface{}, error) {

	result, err := sshClient.RunCommand(diskCommand)
	if err != nil {
		return nil, err
	}

	// Split the output into lines
	lines := strings.Split(strings.TrimSpace(result), "\n")

	disks := make([]map[string]interface{}, 0)

	// Process each line of the command output
	for _, line := range lines {
		fields := strings.Fields(line)

		disk := map[string]interface{}{
			"system.disk.volume":                utils.ParseResult(fields[0]),
			"system.disk.volume.mount.path":     utils.ParseResult(fields[1]),
			"system.disk.volume.used.percent":   utils.ParseResult(strings.TrimSuffix(fields[2], "%")),
			"system.disk.volume.capacity.bytes": utils.ParseResult(fields[3]),
			"system.disk.volume.used.bytes":     utils.ParseResult(fields[4]),
			"system.disk.volume.free.bytes":     utils.ParseResult(fields[5]),
			"system.disk.volume.free.percent":   100 - utils.ParseResult(strings.TrimSuffix(fields[2], "%")).(float64),
		}

		disks = append(disks, disk)
	}

	return disks, nil
}
