package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	client *ssh.Client
}

func ConnectSSH(ip, userName, userPassword string) (*SSHClient, error) {
	config := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.Password(userPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", ip), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH server: %v", err)
	}
	return &SSHClient{client: client}, nil
}

func (s *SSHClient) RunCommand(cmd string) (string, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.Output(cmd)
	return string(output), err
}
