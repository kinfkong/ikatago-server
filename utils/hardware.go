package utils

import (
	"runtime"

	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/pci"
	"github.com/jaypipes/pcidb"
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
	gpu, err := ghw.GPU()
	if err != nil {
		// log.Warnf("Error getting GPU info: %v", err)
		// log.Warnf("Error getting GPU info: %v", err)
		if runtime.GOOS == "darwin" {
			// mock for mac to make it work
			hardwareInfo.GPU = &ghw.GPUInfo{
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
			}
		}
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
