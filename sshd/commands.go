package sshd

import (
	"errors"
	"log"
	"os/exec"

	"github.com/jessevdk/go-flags"
	"github.com/kinfkong/ikatago-server/katago"
)

func init() {
	RegisterCommandHandler("run-katago", runKatago)
}

var runKatagoOpts struct {
	name   *string `long:"name" description:"the katago bin name"`
	weight *string `long:"weight" description:"the katago weight name"`
	config *string `long:"config" description:"the katago config name"`
}

func runKatago(args ...string) (*exec.Cmd, error) {
	_, err := flags.ParseArgs(&runKatagoOpts, args)
	if err != nil {
		log.Printf("ERROR failed to parse kagato args: %+v\n", args)
		return nil, errors.New("invalid_command_args")
	}
	return katago.GetManager().Run(runKatagoOpts.name, runKatagoOpts.weight, runKatagoOpts.config)
}
