package sshd

import (
	"encoding/json"
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
	"github.com/kinfkong/ikatago-server/model"
	"github.com/kinfkong/ikatago-server/utils"
)

func init() {
	RegisterCommandHandler("run-katago", runKatago)
	RegisterCommandHandler("query-server", queryServer)
	RegisterCommandHandler("scp-config", copyConfig)

}

type runkatagoOptsType struct {
	Name            *string `long:"name" description:"the katago bin name"`
	Weight          *string `long:"weight" description:"the katago weight name"`
	Config          *string `long:"config" description:"the katago config name"`
	CustomConfig    *string `long:"custom-config" description:"the katago custom config file name"`
	Compress        bool    `long:"compress" description:"compress the data during transmission"`
	RefreshInterval int     `long:"refresh-interval" description:"sets the refresh interval in cent seconds" default:"30"`
	TransmitMoveNum int     `long:"transmit-move-num" description:"limits number of moves when transmission during analyze" default:"20"`
}

func isValidArg(cmd string, arg string) bool {
	knownCmds := map[string][]string{
		"gtp":       {"-model", "-config"},
		"benchmark": {"-model", "-config"},
		"genconfig": {"-model"},
		"analysis":  {"-model", "-config"},
		"tuner":     {"-model"},
	}
	args, ok := knownCmds[cmd]
	if !ok {
		return false
	}
	for _, v := range args {
		if v == arg {
			return true
		}
	}
	return false

}
func runKatago(session ssh.Session, args ...string) (*exec.Cmd, error) {
	runKatagoOpts := runkatagoOptsType{}
	subcommands, err := flags.ParseArgs(&runKatagoOpts, args)
	if err != nil {
		log.Printf("ERROR failed to parse kagato args: %+v\n", args)
		return nil, errors.New("invalid_command_args")
	}
	if len(subcommands) == 0 {
		// gtp by default
		subcommands = append(subcommands, "gtp")
	}
	found := false
	for _, subcommand := range subcommands {
		if subcommand == "-model" {
			found = true
			break
		}
	}
	if !found {
		if isValidArg(subcommands[0], "-model") {
			subcommands = append(subcommands, "-model", "KATA_WEIGHT_PLACEHOLDER")
		}
	}
	found = false
	for _, subcommand := range subcommands {
		if subcommand == "-config" {
			found = true
			break
		}
	}
	if !found {
		if isValidArg(subcommands[0], "-config") {
			subcommands = append(subcommands, "-config", "KATA_CONFIG_PLACEHOLDER")
		}
	}

	subcommands, err = replaceKataConfigPlaceHolder(session, &runKatagoOpts, subcommands)
	if err != nil {
		return nil, err
	}
	subcommands, err = replaceKataWeightPlaceHolder(session, &runKatagoOpts, subcommands)
	if err != nil {
		return nil, err
	}
	binName := katago.GetManager().DefaultBinName
	if runKatagoOpts.Name != nil {
		binName = *runKatagoOpts.Name
	}
	cmd, err := katago.GetManager().Run(binName, subcommands)
	if err != nil {
		return nil, err
	}
	cmd.Env = session.Environ()
	cmd.Stdin = session
	cmd.Stderr = session.Stderr()

	reader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	gtpWriter := katago.NewGTPWriter(session)
	gtpWriter.MinRefreshCentSecond = runKatagoOpts.RefreshInterval
	gtpWriter.Compression = runKatagoOpts.Compress
	gtpWriter.NumOfTransmitMoves = runKatagoOpts.TransmitMoveNum
	go func() {
		defer reader.Close()
		buffer := make([]byte, 1024*10)
		for {
			n, err := reader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					// done
					// log.Printf("DEBUG end reading")
					break
				} else {
					log.Printf("ERROR failed read buffer: %+v\n", err)
					break
				}
			}
			// write to the session
			gtpWriter.Write(buffer[:n])
		}
		session.Exit(0)
	}()
	return cmd, nil
}

func replaceKataConfigPlaceHolder(session ssh.Session, runKatagoOpts *runkatagoOptsType, subcommands []string) ([]string, error) {
	m := katago.GetManager()
	var configFile *string = nil
	if runKatagoOpts.CustomConfig != nil {
		theFile := fmt.Sprintf("%s/%s/%s", m.CustomConfigDir, session.User(), *runKatagoOpts.CustomConfig)
		configFile = &theFile
	}
	if configFile == nil {
		// no custom config file, use the built-in configs
		configName := runKatagoOpts.Config
		if configName == nil {
			configName = &m.DefaultConfigName
		}
		for _, item := range m.Configs {
			if item.Name == *configName {
				configFile = &item.Path
				break
			}
		}
	}
	result := make([]string, len(subcommands))
	for i, command := range subcommands {
		if command == "KATA_CONFIG_PLACEHOLDER" {
			if configFile == nil {
				return nil, errors.New("no_config_file")
			}
			result[i] = *configFile
		} else {
			result[i] = command
		}
	}
	return result, nil
}

func replaceKataWeightPlaceHolder(session ssh.Session, runKatagoOpts *runkatagoOptsType, subcommands []string) ([]string, error) {
	m := katago.GetManager()
	// no custom weight file, use the built-in weight
	weightName := runKatagoOpts.Weight
	if weightName == nil {
		weightName = &m.DefaultWeightName
	}
	var weightFile *string = nil
	for _, item := range m.Weights {
		if item.Name == *weightName {
			weightFile = &item.Path
			break
		}
	}
	result := make([]string, len(subcommands))
	for i, command := range subcommands {
		if command == "KATA_WEIGHT_PLACEHOLDER" {
			if weightFile == nil {
				return nil, errors.New("no_weight_file")
			}
			result[i] = *weightFile
		} else {
			result[i] = command
		}
	}
	return result, nil
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

func queryServer(session ssh.Session, args ...string) (*exec.Cmd, error) {
	katagoManager := katago.GetManager()
	kataNames := make([]model.ServerInfoItem, 0)
	kataConfigs := make([]model.ServerInfoItem, 0)
	kataWeights := make([]model.ServerInfoItem, 0)
	for _, bin := range katagoManager.Bins {
		kataNames = append(kataNames, model.ServerInfoItem{
			Name:        bin.Name,
			Description: bin.Description,
		})
	}
	for _, weight := range katagoManager.Weights {
		kataWeights = append(kataWeights, model.ServerInfoItem{
			Name:        weight.Name,
			Description: weight.Description,
		})
	}
	for _, configItem := range katagoManager.Configs {
		kataConfigs = append(kataConfigs, model.ServerInfoItem{
			Name:        configItem.Name,
			Description: configItem.Description,
		})
	}
	serverInfo := model.ServerInfo{
		ServerVersion:      utils.ServerVersion,
		SupportKataNames:   kataNames,
		SupportKataConfigs: kataConfigs,
		SupportKataWeights: kataWeights,
		DefaultKataWeight:  katagoManager.DefaultWeightName,
		DefaultKataConfig:  katagoManager.DefaultConfigName,
		DefaultKataName:    katagoManager.DefaultBinName,
	}
	b, err := json.Marshal(serverInfo)
	if err != nil {
		log.Printf("ERROR failed to json the server info: %v\n", err)
		return nil, err
	}
	io.WriteString(session, string(b))
	return nil, nil
}

func copyConfig(session ssh.Session, args ...string) (*exec.Cmd, error) {
	if len(args) == 0 {
		return nil, errors.New("no_target_file_name")
	}
	katagoManager := katago.GetManager()
	outputDir := fmt.Sprintf("%s/%s", katagoManager.CustomConfigDir, session.User())
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, 0755)
	}
	outputFile := fmt.Sprintf("%s/%s", outputDir, args[0])

	buf := new(strings.Builder)
	_, err := io.Copy(buf, session)
	if err != nil {
		log.Printf("ERROR failed to read session: %+v\n", err)
		return nil, err
	}
	err = ioutil.WriteFile(outputFile, []byte(buf.String()), 0644)
	if err != nil {
		log.Printf("ERROR failed to write file: %+v\n", err)
		return nil, err
	}
	return nil, nil
}
