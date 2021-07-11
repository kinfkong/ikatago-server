// +build windows

package utils

import (
	"os/exec"
	"strconv"
)

func KillProcess(pid int) error {
	killCmd := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(pid))
	return killCmd.Run()
}

func SetSysAttr(cmd *exec.Cmd) {
	// does nothing
}
