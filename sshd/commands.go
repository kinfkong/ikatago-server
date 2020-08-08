package sshd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/jessevdk/go-flags"
	"github.com/kinfkong/ikatago-server/katago"
)

func init() {
	RegisterCommandHandler("run-katago", runKatago)
	RegisterCommandHandler("query-katago", queryKatago)
	RegisterCommandHandler("scp-config", copyConfig)

}

var runKatagoOpts struct {
	Name         *string `long:"name" description:"the katago bin name"`
	Weight       *string `long:"weight" description:"the katago weight name"`
	Config       *string `long:"config" description:"the katago config name"`
	CustomConfig *string `long:"custom-config" description:"the katago custom config file name"`
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
	var customConfigFile *string = nil
	if runKatagoOpts.CustomConfig == nil {
		io.WriteString(session, fmt.Sprintf("using katago config: %s\n", configName))
	} else {
		io.WriteString(session, fmt.Sprintf("using custom katago config: %s\n", *runKatagoOpts.CustomConfig))
		// construct the file path
		theFile := fmt.Sprintf("%s/%s/%s", katagoManager.CustomConfigDir, session.User(), *runKatagoOpts.CustomConfig)
		customConfigFile = &theFile
	}

	return katago.GetManager().Run(binName, weightName, configName, customConfigFile)
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

func copyConfig(session ssh.Session, args ...string) (*exec.Cmd, error) {
	if len(args) == 0 {
		return nil, errors.New("no_target_file_name")
	}
	katagoManager := katago.GetManager()
	outputDir := fmt.Sprintf("%s/%s", katagoManager.CustomConfigDir, session.User())
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}
	outputFile := fmt.Sprintf("%s/%s", outputDir, args[0])

	buf := new(strings.Builder)
	_, err := io.Copy(buf, session)
	if err != nil {
		log.Printf("ERROR failed to read session: %+v\n", err)
		return nil, err
	}
	// fmt.Println(buf.String())

	/*f, err := os.Create(outputFile)
	defer f.Close()
	_, err = io.Copy(f, session)
	if err != nil {
		log.Printf("ERROR failed to read session: %+v\n", err)
		return nil, err
	}*/
	err = ioutil.WriteFile(outputFile, []byte(buf.String()), 0644)
	if err != nil {
		log.Printf("ERROR failed to write file: %+v\n", err)
		return nil, err
	}
	return exec.Command("echo", "Copy Done!"), nil
}
