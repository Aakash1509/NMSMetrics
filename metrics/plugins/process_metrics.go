package plugins

import (
	"Plugin/metrics/clients"
	"Plugin/metrics/utils"
)

var processCommands = map[string]string{
	"system.process.memory.used.percent": "ps -eo pmem | awk '{sum+=$1} END {print sum}'",
	"system.process.cpu.percent":         "ps -eo pcpu | awk '{sum+=$1} END {print sum}'",
	"system.process.threads":             "ps -eLf | wc -l",
}

func GetProcessMetrics(sshClient *clients.SSHClient) (map[string]interface{}, error) {

	results := make(map[string]interface{})

	for key, cmd := range processCommands {
		result, err := sshClient.RunCommand(cmd)
		if err != nil {
			results[key] = "ERROR"
			continue // Skip errors
		}
		results[key] = utils.ParseResult(result)
	}

	return results, nil
}
