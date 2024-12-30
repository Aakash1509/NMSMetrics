package clients

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

type SSHClient struct {
	Client *ssh.Client
}

func ConnectSSH(ip string, port int, userName, userPassword string) (*SSHClient, error) {
	config := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.Password(userPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         60 * time.Second,
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH server: %v", err)
	}

	return &SSHClient{Client: conn}, nil
}

func (s *SSHClient) RunCommand(cmd string) (string, error) {
	session, err := s.Client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.Output(cmd)
	return string(output), err
}
