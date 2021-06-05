package cluster

import "time"

type ReportData struct {
	Endpoint   string    `json:"endpoint"`
	MachineID  string    `json:"machineId"`
	ReportType string    `json:"reportType"`
	ReportTime time.Time `json:"reportTime"`
	Data       string    `json:"data"` // json + base64 encoded string
}
