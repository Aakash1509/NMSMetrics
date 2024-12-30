package plugins

import (
	"Plugin/metrics/clients"
	"Plugin/metrics/utils"
)

var deviceCommands = map[string]string{
	"system.network.in.bytes.rate": "cat /proc/net/dev | awk '/lo:/ {print $2}'",
	"system.load.avg1.min":         "uptime | awk -F'load average:' '{ print $2 }' | cut -d',' -f1",
	"system.load.avg5.min":         "uptime | awk -F'load average:' '{ print $2 }' | cut -d',' -f2",
	"system.load.avg15.min":        "uptime | awk -F'load average:' '{ print $2 }' | cut -d',' -f3",
	"system.vendor":                "cat /sys/devices/virtual/dmi/id/sys_vendor",
	"system.os.name":               "uname -o",
	"system.cpu.cores":             "nproc",
	"system.model":                 "cat /sys/devices/virtual/dmi/id/product_name",
	"system.running.processes":     "ps -e | wc -l",
	"system.blocked.processes":     "grep \"procs_blocked\" /proc/stat | awk '{print $2}'",
}

func GetDeviceMetrics(sshClient *clients.SSHClient) (map[string]interface{}, error) {

	results := make(map[string]interface{})

	for key, cmd := range deviceCommands {
		result, err := sshClient.RunCommand(cmd)
		if err != nil {
			results[key] = "ERROR"
			continue // Skip errors
		}
		results[key] = utils.ParseResult(result)
	}

	return results, nil
}
