package event

import (
	"sync"
)

// Service type
type Service struct {
	bus Bus
}

var serviceInstance *Service
var serviceMu sync.Mutex

// GetService returns the singleton instance of the Service
func GetService() *Service {
	serviceMu.Lock()
	defer serviceMu.Unlock()
	if serviceInstance == nil {
		bus := New()
		serviceInstance = &Service{
			bus: bus,
		}
	}

	return serviceInstance
}

// Publish publishes the event
func (service *Service) Publish(topic string, args ...interface{}) {
	service.bus.Publish(topic, args...)
}

// Subscribe subscribes the event
func (service *Service) Subscribe(topic string, fn interface{}) error {
	return service.bus.Subscribe(topic, fn)
}
