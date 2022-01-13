// +build !windows

package utils

import (
	"log"
	"syscall"

	ps "github.com/mitchellh/go-ps"
)

func KillProcessAndChildren(pid int) error {
	processes, err := ps.Processes()
	if err != nil {
		log.Printf("ERROR failed to get proccesses: %+v", err)
		return err
	}
	return _killProcessAndChildren(pid, processes)
}

func _killProcessAndChildren(pid int, procceses []ps.Process) error {
	// kill the children first
	for _, p := range procceses {
		if p.PPid() == pid {
			err := _killProcessAndChildren(p.Pid(), procceses)
			if err != nil {
				log.Printf("Failed to kill process: %+v", err)
				return err
			}
		}
	}
	// kill the current process
	return syscall.Kill(pid, syscall.SIGKILL)

}
