package utils

import (
	"errors"
	"os/exec"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

var cmdManagerInstance *CmdManager
var cmdManagerMu sync.Mutex

type OnClientClosedHandler func(error)

type ExtendedCmd struct {
	// ID is the id of the command, assigned internally
	ID          string
	CommandType string `json:"commandType"`

	// Cmd the command it self
	Cmd *exec.Cmd
	// Username the ikatago user that runs this command
	Username       *string
	StartedAt      *time.Time
	OnClientClosed OnClientClosedHandler
}

type CommandInfo struct {
	ID          string     `json:"id"`
	CommandType string     `json:"commandType"`
	Username    *string    `json:"username"`
	StartedAt   *time.Time `json:"startedAt"`
	Path        string     `json:"path"`
	Args        []string   `json:"args"`
	Env         []string   `json:"env"`
	Dir         string     `json:"dir"`
	Pid         *int       `json:"pid"`
}

type CmdManager struct {
	cmds []*ExtendedCmd
}

// GetCmdManager returns the singleton instance of the cmd manager
func GetCmdManager() *CmdManager {
	cmdManagerMu.Lock()
	defer cmdManagerMu.Unlock()

	if cmdManagerInstance == nil {
		cmdManagerInstance = &CmdManager{
			cmds: make([]*ExtendedCmd, 0),
		}
	}
	return cmdManagerInstance
}

// RunCommand runs the command sync (will block until the cmd run done)
func (cmdManager *CmdManager) RunCommand(cmd *ExtendedCmd) error {
	if cmd == nil {
		return errors.New("cmd cannot be nil")
	}

	now := time.Now()
	cmd.StartedAt = &now
	cmdManager.addCmd(cmd)

	err := cmd.Cmd.Run()
	// remove it after the command done
	cmdManager.removeCmdByID(cmd.ID)
	return err
}

// PrepareCommand runs the command sync (will block until the cmd run done)
func (cmdManager *CmdManager) PrepareCommand(username *string, commandType string, cmd *exec.Cmd) (*ExtendedCmd, error) {
	if cmd == nil {
		return nil, errors.New("cmd cannot be nil")
	}

	now := time.Now()
	// add to the current commands
	cmdID := uuid.NewV4().String()
	result := &ExtendedCmd{
		ID:          cmdID,
		CommandType: commandType,
		Username:    username,
		Cmd:         cmd,
		StartedAt:   &now,
	}
	// change the cmd stdin
	rawReader := cmd.Stdin
	reader := NewIOReaderWrapper(rawReader)
	reader.onClientClosed = func(err error) {
		if result.OnClientClosed != nil {
			result.OnClientClosed(err)
		}
	}
	cmd.Stdin = reader
	return result, nil
}

func (cmdManager *CmdManager) KillCommand(cmdID string) error {
	for _, cmd := range cmdManager.cmds {
		if cmd.ID == cmdID && cmd.Cmd.Process != nil {
			err := KillProcessAndChildren(cmd.Cmd.Process.Pid)
			if err != nil {
				return err
			}
			break
		}
	}
	// not found
	return nil
}

func (cmdManager *CmdManager) GetAllCmdInfo() []CommandInfo {
	cmdManagerMu.Lock()
	defer cmdManagerMu.Unlock()
	infos := make([]CommandInfo, 0)
	for _, cmd := range cmdManager.cmds {
		var pid *int = nil
		if cmd.Cmd.Process != nil {
			pid = &cmd.Cmd.Process.Pid
		}
		infos = append(infos, CommandInfo{
			ID:          cmd.ID,
			CommandType: cmd.CommandType,
			Username:    cmd.Username,
			StartedAt:   cmd.StartedAt,
			Path:        cmd.Cmd.Path,
			Args:        cmd.Cmd.Args,
			Env:         cmd.Cmd.Env,
			Dir:         cmd.Cmd.Dir,
			Pid:         pid,
		})
	}
	return infos
}

func (cmdManager *CmdManager) addCmd(cmd *ExtendedCmd) {
	cmdManagerMu.Lock()
	defer cmdManagerMu.Unlock()
	cmdManager.cmds = append(cmdManager.cmds, cmd)
}

func (cmdManager *CmdManager) removeCmdByID(ID string) {
	cmdManagerMu.Lock()
	defer cmdManagerMu.Unlock()
	foundIndex := -1
	for i := range cmdManager.cmds {
		if cmdManager.cmds[i].ID == ID {
			foundIndex = i
			break
		}
	}
	if foundIndex < 0 {
		// not found
		return
	}
	// remove that from index
	cmdManager.cmds = append(cmdManager.cmds[:foundIndex], cmdManager.cmds[foundIndex+1:]...)
}
