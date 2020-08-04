package main

import (
	"github.com/kinfkong/ikatago-server/nat"
	"github.com/kinfkong/ikatago-server/sshd"
)

func main() {
	go sshd.Run()
	//  ssh -p 10020 akv100kt6@gpu61.mistgpu.com
	natProvider := &nat.Knat{
		SSHHost:     "120.53.123.43",
		SSHPort:     22,
		SSHUsername: "nat",
		SSHPassword: "IamsureKKNo.1",
		LocalPort:   2222,
	}
	natProvider.Run()
}
