package cluster

import "sync"

// Service type
type Service struct {
	Server    string  `json:"server"`
	MachineID *string `json:"machineId"`
}

var serviceInstance *Service
var serviceMu sync.Mutex

// GetService returns the singleton instance of the Service
func GetService() *Service {
	serviceMu.Lock()
	defer serviceMu.Unlock()

	if serviceInstance == nil {
		serviceInstance = &Service{}
	}
	return serviceInstance
}

func Init(clusterToken string) error {
	return nil
}
