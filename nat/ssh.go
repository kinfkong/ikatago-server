package nat

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHNatProvider represents the ssh nat provider
type SSHNatProvider struct {
	SSHHost string `json:"sshHost"`
	SSHPort int    `json:"sshPort"`

	SSHUsername string `json:"sshUsername"`
	SSHPassword string `json:"sshPassword"`

	RemotePort int `json:"remote_port"`
	LocalPort  int `json:"local_port"`
}

var _ Provider = (&SSHNatProvider{})

// Run runs the nat service
func (nat *SSHNatProvider) Run() error {
	// Connection settings
	sshAddr := fmt.Sprintf("%s:%d", nat.SSHHost, nat.SSHPort)
	localAddr := fmt.Sprintf("%s:%d", "127.0.0.1", nat.LocalPort)
	remoteAddr := fmt.Sprintf("%s:%d", "0.0.0.0", nat.RemotePort)

	// Build SSH client configuration
	cfg, err := makeSSHConfig(nat.SSHUsername, nat.SSHPassword)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	// Establish connection with SSH server
	conn, err := ssh.Dial("tcp", sshAddr, cfg)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	defer conn.Close()
	log.Printf("DEBUG: success to ssh to: %s\n", sshAddr)
	// Listen on remote server port
	listener, err := conn.Listen("tcp", remoteAddr)
	if err != nil {
		log.Fatalln(fmt.Printf("Listen open port ON remote server error: %s", err))
	}
	defer listener.Close()
	log.Printf("DEBUG: success to create listenning on remote %s\n", remoteAddr)

	// handle incoming connections on reverse forwarded tunnel
	for {
		// Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
		local, err := net.Dial("tcp", localAddr)
		if err != nil {
			log.Printf("DEBUG: cannot connected to local: " + localAddr)
			time.Sleep(time.Second * 1)
			continue
		}
		defer local.Close()
		// log.Printf("DEBUG: connected to local: " + localAddr)
		client, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		handleClient(client, local)
	}
}

// GetInfo gets the info
func (nat *SSHNatProvider) GetInfo() (Info, error) {
	return Info{
		Host: nat.SSHHost,
		Port: nat.SSHPort,
	}, nil
}

// Get ssh client config for our connection
// SSH config will use 2 authentication strategies: by key and by password
func makeSSHConfig(username string, password string) (*ssh.ClientConfig, error) {
	config := ssh.ClientConfig{
		Timeout:         30 * time.Second,
		User:            username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	config.Auth = []ssh.AuthMethod{ssh.Password(password)}

	return &config, nil
}

func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	// Start remote -> local data transfer
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			log.Println(fmt.Sprintf("error while copy remote->local: %s", err))
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			log.Println(fmt.Sprintf("error while copy local->remote: %s", err))
		}
		chDone <- true
	}()

	<-chDone
}
