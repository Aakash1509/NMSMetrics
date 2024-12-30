package plugins

import (
	"Plugin/metrics/clients"
	"Plugin/metrics/utils"
)

var cpuCommands = map[string]string{
	"system.cpu.core.idle.percent":      "top -bn1 | grep \"Cpu(s)\" | awk '{print $8}'",
	"system.cpu.core.percent":           "top -bn1 | grep \"Cpu(s)\" | awk '{print $2 + $4 + $6 + $8 + $10}'",
	"system.cpu.core.user.percent":      "top -bn1 | grep \"Cpu(s)\" | awk '{print $2}'",
	"system.cpu.core.kernel.percent":    "top -bn1 | grep \"Cpu(s)\" | awk '{print $4}'",
	"system.cpu.core.io.percent":        "top -bn1 | grep \"Cpu(s)\" | awk '{print $6}'",
	"system.cpu.core.interrupt.percent": "top -bn1 | grep \"Cpu(s)\" | awk '{print $12 + $14}'",
}

func GetCPUMetrics(sshClient *clients.SSHClient) (map[string]interface{}, error) {

	results := make(map[string]interface{})

	for key, cmd := range cpuCommands {
		result, err := sshClient.RunCommand(cmd)
		if err != nil {
			results[key] = "ERROR"
			continue // Skip errors
		}
		results[key] = utils.ParseResult(result)
	}

	return results, nil
}
