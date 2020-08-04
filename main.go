package main

import (
	"github.com/kinfkong/ikatago-server/nat"
	"github.com/kinfkong/ikatago-server/sshd"
)

func main() {
	go sshd.Run()
	natProvider := &nat.Knat{
		SSHHost:     "120.53.123.43",
		SSHPort:     8203,
		SSHUsername: "nat",
		SSHPassword: "IamsureKKNo.1",
		LocalPort:   2222,
	}
	natProvider.Run()
}
