package nat

import (
	"errors"
	"log"
	"strconv"
)

// Direct means the direct connect, don't use any nat
type Direct struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

var _ Provider = (&Direct{})

// InitWithConfig inits from the config
func (direct *Direct) InitWithConfig(configObject map[string]interface{}) error {
	hostV, ok := configObject["host"]
	if !ok {
		log.Printf("ERROR missing host config for the direct nat.")
		return errors.New("missing_host")
	}
	host, ok := hostV.(string)
	if !ok {
		log.Printf("ERROR host configure should be a string")
		return errors.New("invalid_host")
	}
	portV, ok := configObject["port"]
	if !ok {
		log.Printf("ERROR missing port config for the direct nat.")
		return errors.New("missing_port")
	}
	port, ok := portV.(int)
	var err error
	if !ok {
		port, err = strconv.Atoi(portV.(string))
	}
	if err != nil {
		log.Printf("ERROR cannot convert port to a number")
		return errors.New("invalid_port")
	}
	direct.Host = host
	direct.Port = port
	return nil
}

// RunAsync means runs async
func (direct *Direct) RunAsync() error {
	// does nothing
	return nil
}

// GetInfo gets the info
func (direct *Direct) GetInfo() (Info, error) {
	return Info{
		Host: direct.Host,
		Port: direct.Port,
	}, nil
}
