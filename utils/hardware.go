package utils

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/pci"
	"github.com/jaypipes/pcidb"
	"github.com/kinfkong/ikatago-server/errors"
)

type HardwareInfo struct {
	GPU    *ghw.GPUInfo    `json:"gpu"`
	CPU    *ghw.CPUInfo    `json:"cpu"`
	Block  *ghw.BlockInfo  `json:"block"`
	Memory *ghw.MemoryInfo `json:"memory"`
}
type GPUInfo struct {
	GPUs []string `json:"gpus"`
}

var gHardwareInfo *HardwareInfo = nil
var gGPUInfo *GPUInfo = nil

func GetGPUInfo() []string {
	if gGPUInfo != nil {
		return gGPUInfo.GPUs
	}
	gpuInfo := &GPUInfo{
		GPUs: make([]string, 0),
	}
	hardwareInfo := GetHardwareInfo()
	if hardwareInfo == nil || hardwareInfo.GPU == nil {
		// log.Warnf("Error getting GPU info: %v", err)
	} else {
		if hardwareInfo.GPU.GraphicsCards != nil {
			for _, card := range hardwareInfo.GPU.GraphicsCards {
				if card == nil || card.DeviceInfo == nil || card.DeviceInfo.Product == nil {
					continue
				}
				gpuInfo.GPUs = append(gpuInfo.GPUs, card.DeviceInfo.Product.Name)
			}
		}
	}

	gGPUInfo = gpuInfo
	return gpuInfo.GPUs
}

func GetHardwareInfo() *HardwareInfo {
	if gHardwareInfo != nil {
		return gHardwareInfo
	}
	hardwareInfo := &HardwareInfo{}
	gpu, err := interceptGPUInfo()
	if gpu == nil {
		gpu, err = ghw.GPU()
	}
	if err != nil {
	} else {
		hardwareInfo.GPU = gpu
	}
	cpu, err := ghw.CPU()
	if err != nil {
		// log.Warnf("Error getting CPU info: %v", err)
	} else {
		hardwareInfo.CPU = cpu
	}

	block, err := ghw.Block()
	if err != nil {
		// log.Warnf("Error getting block storage info: %v", err)
	} else {
		hardwareInfo.Block = block
	}

	memory, err := ghw.Memory()
	if err != nil {
		// log.Warnf("Error getting memory info: %v", err)
	} else {
		hardwareInfo.Memory = memory
	}
	gHardwareInfo = hardwareInfo
	return hardwareInfo
}

func interceptGPUInfo() (*ghw.GPUInfo, error) {
	if runtime.GOOS == "darwin" {
		// mock for mac to make it work
		return &ghw.GPUInfo{
			GraphicsCards: []*ghw.GraphicsCard{
				{
					DeviceInfo: &pci.Device{
						Product: &pcidb.Product{
							Name: "Apple M1 2020",
						},
						Vendor:               &pcidb.Vendor{},
						Subsystem:            &pcidb.Product{},
						Class:                &pcidb.Class{},
						Subclass:             &pcidb.Subclass{},
						ProgrammingInterface: &pcidb.ProgrammingInterface{},
					},
				},
			},
		}, nil
	} else if runtime.GOOS == "linux" {
		// read from nvidia-smi
		cmd := exec.Command("bash", "-c", "nvidia-smi -q | grep \"Product Name\" | sed \"s/.*Product Name.*:\\(.*\\)$/\\1/g\"")
		response, err := cmd.CombinedOutput()
		if err != nil {
			return nil, errors.CreateError(400, "nvidia not supported")
		}
		gpuInfo := &ghw.GPUInfo{
			GraphicsCards: []*ghw.GraphicsCard{},
		}
		// parse the response
		products := strings.Split(string(response), "\n")
		for _, productName := range products {
			if len(productName) == 0 {
				continue
			}
			name := strings.TrimSpace(productName)
			gpuInfo.GraphicsCards = append(gpuInfo.GraphicsCards, &ghw.GraphicsCard{
				DeviceInfo: &pci.Device{
					Product: &pcidb.Product{
						Name: name,
					},
					Vendor:               &pcidb.Vendor{},
					Subsystem:            &pcidb.Product{},
					Class:                &pcidb.Class{},
					Subclass:             &pcidb.Subclass{},
					ProgrammingInterface: &pcidb.ProgrammingInterface{},
				},
			})
		}
		return gpuInfo, nil
	}
	return nil, nil
}
