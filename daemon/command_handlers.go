package daemon

import (
	"log"

	"github.com/kinfkong/ikatago-server/errors"
	"github.com/kinfkong/ikatago-server/utils"
)

func KillCommandHandler(command *ResponseCommand) error {
	log.Printf("INFO got kill command")
	if command == nil || command.Command != "kill" || len(command.Args) < 1 {
		return errors.CreateError(400, "invalid_command")
	}
	cmdID := command.Args[0]

	log.Printf("INFO killing command: %v", cmdID)
	return utils.GetCmdManager().KillCommand(cmdID)
}
