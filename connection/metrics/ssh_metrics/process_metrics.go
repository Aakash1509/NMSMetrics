package ssh_metrics

import (
	"Plugin/connection/clients"
	"Plugin/connection/utils/parse"
)

type ProcessMetrics struct {
	MemoryUsedPercent interface{} `json:"system.process.memory.used.percent"`
	CPUProcessPercent interface{} `json:"system.process.cpu.percent"`
	Threads           interface{} `json:"system.process.threads"`
}

func GetProcessMetrics(sshClient *clients.SSHClient) (*ProcessMetrics, error) {
	metrics := &ProcessMetrics{}

	commands := map[string]*interface{}{
		"ps -eo pmem | awk '{sum+=$1} END {print sum}'": &metrics.MemoryUsedPercent,
		"ps -eo pcpu | awk '{sum+=$1} END {print sum}'": &metrics.CPUProcessPercent,
		"ps -eLf | wc -l": &metrics.Threads,
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
