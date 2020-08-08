package sshd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/jessevdk/go-flags"
	"github.com/kinfkong/ikatago-server/katago"
)

func init() {
	RegisterCommandHandler("run-katago", runKatago)
	RegisterCommandHandler("query-katago", queryKatago)

}

var runKatagoOpts struct {
	Name   *string `long:"name" description:"the katago bin name"`
	Weight *string `long:"weight" description:"the katago weight name"`
	Config *string `long:"config" description:"the katago config name"`
}

func runKatago(session ssh.Session, args ...string) (*exec.Cmd, error) {
	_, err := flags.ParseArgs(&runKatagoOpts, args)
	if err != nil {
		log.Printf("ERROR failed to parse kagato args: %+v\n", args)
		return nil, errors.New("invalid_command_args")
	}
	outputKataInfo(session)
	katagoManager := katago.GetManager()

	binName, weightName, configName := katagoManager.GetCurrentUsingNames(runKatagoOpts.Name, runKatagoOpts.Weight, runKatagoOpts.Config)
	io.WriteString(session, fmt.Sprintf("using katago name: %s\n", binName))
	io.WriteString(session, fmt.Sprintf("using katago weight: %s\n", weightName))
	io.WriteString(session, fmt.Sprintf("using katago config: %s\n", configName))
	return katago.GetManager().Run(binName, weightName, configName)
}

func outputKataInfo(session ssh.Session) {
	katagoManager := katago.GetManager()
	weights := make([]string, 0)
	for _, weight := range katagoManager.Weights {
		weights = append(weights, weight.Name)
	}
	bins := make([]string, 0)
	for _, bin := range katagoManager.Bins {
		bins = append(bins, bin.Name)
	}
	configs := make([]string, 0)
	for _, kataConfig := range katagoManager.Configs {
		configs = append(configs, kataConfig.Name)
	}
	io.WriteString(session, fmt.Sprintf("support katago names: %s\n", strings.Join(bins, ", ")))
	io.WriteString(session, fmt.Sprintf("support katago weights: %s\n", strings.Join(weights, ", ")))
	io.WriteString(session, fmt.Sprintf("support katago configs: %s\n", strings.Join(configs, ", ")))
}

func queryKatago(session ssh.Session, args ...string) (*exec.Cmd, error) {
	outputKataInfo(session)
	katagoManager := katago.GetManager()
	io.WriteString(session, fmt.Sprintf("default katago name: %s\n", katagoManager.DefaultBinName))
	io.WriteString(session, fmt.Sprintf("default katago weight: %s\n", katagoManager.DefaultWeightName))
	io.WriteString(session, fmt.Sprintf("default katago config: %s\n", katagoManager.DefaultConfigName))
	return nil, nil
}
