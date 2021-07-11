// +build !windows

package utils

import (
	"os/exec"
	"syscall"
)

func KillProcess(pid int) error {
	return syscall.Kill(-pid, syscall.SIGKILL)
}

func SetSysAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}
