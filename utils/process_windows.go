// +build windows

package utils

import (
	"os/exec"
	"strconv"
)

func KillProcessAndChildren(pid int) error {
	killCmd := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(pid))
	return killCmd.Run()
}
