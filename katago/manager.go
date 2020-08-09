package katago

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"

	"github.com/kinfkong/ikatago-server/config"
	"github.com/kinfkong/ikatago-server/utils"

	"github.com/spf13/viper"
)

// BinConfig the bin configs
type BinConfig struct {
	Name   string  `json:"name"`
	Path   string  `json:"path"`
	Runner *string `json:"runner"`
}

// WeightConfig the bin configs
type WeightConfig struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// ConfigConfig the bin configs
type ConfigConfig struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// Manager managers the katagos
type Manager struct {
	Bins              []BinConfig    `json:"bins"`
	Weights           []WeightConfig `json:"weights"`
	Configs           []ConfigConfig `json:"configs"`
	DefaultBinName    string         `json:"defaultBinName"`
	DefaultWeightName string         `json:"defaultWeightName"`
	DefaultConfigName string         `json:"defaultConfigName"`
	CustomConfigDir   string         `json:"customConfigDir"`
}

var managerInstance *Manager
var managerMu sync.Mutex

// GetManager returns the singleton instance of the Service
func GetManager() *Manager {
	managerMu.Lock()
	defer managerMu.Unlock()

	if managerInstance == nil {
		managerInstance = NewManager(config.GetConfig().Sub("katago"))
	}
	return managerInstance
}

// NewManager creates the kata manager
func NewManager(configObject *viper.Viper) *Manager {
	manager := Manager{}
	err := configObject.Unmarshal(&manager)
	if err != nil {
		log.Printf("Failed to unmarshal config. %+v\n", err)
		return nil
	}
	// validate paths
	for _, bin := range manager.Bins {
		if bin.Runner != nil && *bin.Runner == "aistudio-runner" {
			// special the path is a directly
			if !utils.DirectoryExists(bin.Path) {
				log.Printf("ERROR the path %s does not exist or not a directory\n", bin.Path)
				return nil
			}
		} else {
			if !utils.FileExists(bin.Path) {
				log.Printf("ERROR the path %s does not exist or not a file\n", bin.Path)
				return nil
			}
		}
	}
	for _, weight := range manager.Weights {
		if !utils.FileExists(weight.Path) {
			log.Printf("ERROR the path %s does not exist or not a file\n", weight.Path)
			return nil
		}
	}
	for _, config := range manager.Configs {
		if !utils.FileExists(config.Path) {
			log.Printf("ERROR the path %s does not exist or not a file\n", config.Path)
			return nil
		}
	}
	if len(manager.Weights) == 0 {
		log.Printf("ERROR no model weights configured in this server")
		return nil
	}
	if len(manager.Bins) == 0 {
		log.Printf("ERROR no katago binaries configured in this server")
		return nil
	}
	if len(manager.Configs) == 0 {
		log.Printf("ERROR no katago config files configured in this server")
		return nil
	}
	if len(manager.DefaultBinName) == 0 {
		manager.DefaultBinName = manager.Bins[0].Name
	}
	if len(manager.DefaultWeightName) == 0 {
		manager.DefaultWeightName = manager.Weights[0].Name
	}
	if len(manager.DefaultConfigName) == 0 {
		manager.DefaultConfigName = manager.Configs[0].Name
	}
	return &manager
}

func (m *Manager) runDirectly(binPath string, subcommands []string) (*exec.Cmd, error) {
	return exec.Command(binPath, subcommands...), nil
}

func (m *Manager) runByRunner(runnerPath string, binPath string, subcommands []string) (*exec.Cmd, error) {
	all := []string{binPath}
	all = append(all, subcommands...)
	return exec.Command(runnerPath, all...), nil
}

func (m *Manager) runByAiStudioRunner(binName string, binPath string, subcommands []string) (*exec.Cmd, error) {
	decryptePassword := "abcde12345"
	decrypteCommandTemplate := "openssl enc -in %s -d -aes-256-cbc -pass pass:%s > %s"

	inputRootPath := binPath
	outputRootPath := "/tmp/" + binName
	output, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -rf %s && mkdir -p %s", outputRootPath, outputRootPath)).CombinedOutput()
	if err != nil {
		return nil, err
	}
	katagoInputName := fmt.Sprintf("%s/k", inputRootPath)
	katagoOutputName := fmt.Sprintf("%s/k", outputRootPath)
	libzipInputName := fmt.Sprintf("%s/lz", inputRootPath)
	libzipOutputName := fmt.Sprintf("%s/lz", outputRootPath)
	libstdcInputName := fmt.Sprintf("%s/lc", inputRootPath)
	libstdcOutputName := fmt.Sprintf("%s/lc", outputRootPath)

	libzipRealName := fmt.Sprintf("%s/libzip.so.4", outputRootPath)
	libstdcRealName := fmt.Sprintf("%s/libstdc++.so.6", outputRootPath)

	output, err = exec.Command("/bin/sh", "-c", fmt.Sprintf(decrypteCommandTemplate, katagoInputName, decryptePassword, katagoOutputName)).CombinedOutput()
	if err != nil {
		log.Printf("DEBUG error output: %s\n", string(output))
		return nil, err
	}
	output, err = exec.Command("/bin/sh", "-c", fmt.Sprintf(decrypteCommandTemplate, libzipInputName, decryptePassword, libzipOutputName)).CombinedOutput()
	if err != nil {
		log.Printf("DEBUG error output: %s\n", string(output))
		return nil, err
	}
	output, err = exec.Command("/bin/sh", "-c", fmt.Sprintf(decrypteCommandTemplate, libstdcInputName, decryptePassword, libstdcOutputName)).CombinedOutput()
	if err != nil {
		log.Printf("DEBUG error output: %s\n", string(output))
		return nil, err
	}

	output, err = exec.Command("/bin/sh", "-c", fmt.Sprintf("chmod +x %s", katagoOutputName)).CombinedOutput()
	if err != nil {
		log.Printf("DEBUG error output: %s\n", string(output))
		return nil, err
	}
	output, err = exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -f %s && ln -s %s %s", libstdcRealName, libstdcOutputName, libstdcRealName)).CombinedOutput()
	if err != nil {
		log.Printf("DEBUG error output: %s\n", string(output))
		return nil, err
	}

	output, err = exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -f %s && ln -s %s %s", libzipRealName, libzipOutputName, libzipRealName)).CombinedOutput()
	if err != nil {
		log.Printf("DEBUG error output: %s\n", string(output))
		return nil, err
	}

	return exec.Command("/bin/sh", "-c", fmt.Sprintf("export LD_LIBRARY_PATH=%s:$LD_LIBRARY_PATH; %s %s", outputRootPath, katagoOutputName, strings.Join(subcommands, " "))), nil
}

// Run runs the katago
func (m *Manager) Run(binName string, subcommands []string) (*exec.Cmd, error) {

	var binConfig *BinConfig = nil

	for _, item := range m.Bins {
		if item.Name == binName {
			binConfig = &item
			break
		}
	}

	if binConfig == nil {
		log.Printf("bin name: " + binName + " not found.")
		return nil, errors.New("bin name: " + binName + " not found.")
	}

	if binConfig.Runner == nil || len(*binConfig.Runner) == 0 {
		// no runner, run directly
		return m.runDirectly(binConfig.Path, subcommands)
	}
	// run by runner
	if *binConfig.Runner == "aistudio-runner" {
		// special for aistudio
		return m.runByAiStudioRunner(binName, binConfig.Path, subcommands)
	}
	return m.runByRunner(*binConfig.Runner, binConfig.Path, subcommands)
}
