package ssh_metrics

import (
	"Plugin/connection/clients"
	"Plugin/connection/utils/parse"
)

type DeviceMetrics struct {
	NetworkInBytesRate interface{} `json:"system.network.in.bytes.rate"`
	LoadAvg1Min        interface{} `json:"system.load.avg1.min"`
	LoadAvg5Min        interface{} `json:"system.load.avg5.min"`
	LoadAvg15Min       interface{} `json:"system.load.avg15.min"`
	Vendor             interface{} `json:"system.vendor"`
	OSName             interface{} `json:"system.os.name"`
	CPUCores           interface{} `json:"system.cpu.cores"`
	Model              interface{} `json:"system.model"`
	RunningProcesses   interface{} `json:"system.running.processes"`
	BlockedProcesses   interface{} `json:"system.blocked.processes"`
}

func GetDeviceMetrics(sshClient *clients.SSHClient) (*DeviceMetrics, error) {
	metrics := &DeviceMetrics{}

	commands := map[string]*interface{}{
		"cat /proc/net/dev | awk '/lo:/ {print $2}'":                    &metrics.NetworkInBytesRate,
		"uptime | awk -F'load average:' '{ print $2 }' | cut -d',' -f1": &metrics.LoadAvg1Min,
		"uptime | awk -F'load average:' '{ print $2 }' | cut -d',' -f2": &metrics.LoadAvg5Min,
		"uptime | awk -F'load average:' '{ print $2 }' | cut -d',' -f3": &metrics.LoadAvg15Min,
		"cat /sys/devices/virtual/dmi/id/sys_vendor":                    &metrics.Vendor,
		"uname -o": &metrics.OSName,
		"nproc":    &metrics.CPUCores,
		"cat /sys/devices/virtual/dmi/id/product_name": &metrics.Model,
		"ps -e | wc -l": &metrics.RunningProcesses,
		"grep \"procs_blocked\" /proc/stat | awk '{print $2}'": &metrics.BlockedProcesses,
	}

	for cmd, field := range commands {
		result, err := sshClient.RunCommand(cmd)
		if err != nil {
			*field = "ERROR"
			continue // Continue with the next command
		}
		*field = parse.ParseResult(result)
	}
	return metrics, nil
}
