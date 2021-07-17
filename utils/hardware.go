package utils

import (
	"github.com/jaypipes/ghw"
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
	gpu, err := ghw.GPU()
	if err != nil {
		// log.Warnf("Error getting GPU info: %v", err)
	} else {
		if gpu.GraphicsCards != nil {
			for _, card := range gpu.GraphicsCards {
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
