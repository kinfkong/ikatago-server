package nat

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// Knat represents the knat
type Knat struct {
	SSHHost string `json:"sshHost"`
	SSHPort int    `json:"sshPort"`

	SSHUsername string `json:"sshUsername"`
	SSHPassword string `json:"sshPassword"`

	LocalPort int `json:"local_port"`
}

func (knat *Knat) fetchRemotePort() (int, error) {
	config := &ssh.ClientConfig{
		Timeout:         30 * time.Second,
		User:            knat.SSHUsername,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	config.Auth = []ssh.AuthMethod{ssh.Password(knat.SSHPassword)}

	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", knat.SSHHost, knat.SSHPort), config)
	if err != nil {
		log.Fatal("failed to create ssh client", err)
	}
	defer sshClient.Close()

	// start the sesssion to do it
	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatal("failed to create ssh session", err)
	}
	defer session.Close()

	cmd := fmt.Sprintf("/home/nat/assign-port.sh")
	log.Printf("DEBUG running commad:%s\n", cmd)
	output, err := session.Output(cmd)
	if err != nil {
		return 0, err
	}
	log.Printf("DEBUG fetched port: %s\n", string(output))
	remotePort, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, err
	}
	return remotePort, nil
}

// Run runs the nat service
func (knat *Knat) Run() error {
	remotePort, err := knat.fetchRemotePort()
	if err != nil {
		log.Fatal("failed to fetch remote port", err)
	}

	sshProvider := &SSHNatProvider{
		SSHHost:     knat.SSHHost,
		SSHPort:     knat.SSHPort,
		SSHUsername: knat.SSHUsername,
		SSHPassword: knat.SSHPassword,
		RemotePort:  remotePort,
		LocalPort:   knat.LocalPort,
	}
	return sshProvider.Run()
}
