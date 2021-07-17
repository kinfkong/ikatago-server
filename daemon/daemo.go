package daemon

import (
	"time"

	"github.com/kinfkong/ikatago-server/utils"
)

// WorkerData defines the worker data entity
type WorkerData struct {
	// ID the id of thw worker
	ID string `json:"id"`
	// WorkerType the type of the worker: ikatago-server, daemon
	WorkerType string `json:"workerType"`
	// RunningCommands the commands that created by the worker and running
	RunningCommands []utils.CommandInfo `json:"runningCommands"`
	GPUs            []string            `json:"gpus"`
	// ExtraInfo the extraInfo
	ExtraInfo map[string]interface{} `json:"extraInfo"`
	// Timestamp the time stamp of the data
	Timestamp time.Time `json:"timestamp"`
}
