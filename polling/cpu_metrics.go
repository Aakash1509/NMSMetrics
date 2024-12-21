package main

type CPUMetrics struct {
	IdlePercent      interface{} `json:"system.cpu.core.idle.percent"`
	CorePercent      interface{} `json:"system.cpu.core.percent"`
	UserPercent      interface{} `json:"system.cpu.core.user.percent"`
	KernelPercent    interface{} `json:"system.cpu.core.kernel.percent"`
	IOPercent        interface{} `json:"system.cpu.core.io.percent"`
	InterruptPercent interface{} `json:"system.cpu.core.interrupt.percent"`
}

func GetCPUMetrics(sshClient *SSHClient) (*CPUMetrics, error) {
	metrics := &CPUMetrics{}

	commands := map[string]*interface{}{
		"top -bn1 | grep \"Cpu(s)\" | awk '{print $8}'":                      &metrics.IdlePercent,
		"top -bn1 | grep \"Cpu(s)\" | awk '{print $2 + $4 + $6 + $8 + $10}'": &metrics.CorePercent,
		"top -bn1 | grep \"Cpu(s)\" | awk '{print $2}'":                      &metrics.UserPercent,
		"top -bn1 | grep \"Cpu(s)\" | awk '{print $4}'":                      &metrics.KernelPercent,
		"top -bn1 | grep \"Cpu(s)\" | awk '{print $6}'":                      &metrics.IOPercent,
		"top -bn1 | grep \"Cpu(s)\" | awk '{print $12 + $14}'":               &metrics.InterruptPercent,
	}

	for cmd, field := range commands {
		result, err := sshClient.RunCommand(cmd)
		if err != nil {
			*field = "ERROR"
			continue // Continue with the next command
		}
		*field = ParseResult(result)
	}
	return metrics, nil
}
