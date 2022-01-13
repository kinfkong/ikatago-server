package nat

import (
	"errors"
	"log"
	"sync"

	"github.com/kinfkong/ikatago-server/config"
)

var gProvider Provider = nil
var gProviderMutex sync.Mutex

// GetNatProvider gets the nat provider
func GetNatProvider() (Provider, error) {
	gProviderMutex.Lock()
	defer gProviderMutex.Unlock()

	if gProvider != nil {
		return gProvider, nil
	}
	// read the name from config
	natName := config.GetConfig().GetString("use_nat")
	if len(natName) == 0 {
		log.Printf("ERROR nat name is not specified.")
		return nil, errors.New("invalid_use_nat")
	}
	natConfig := config.GetConfig().GetStringMap("nats." + natName)
	if natConfig == nil {
		log.Printf("ERROR cannot find config of: " + "nats." + natName)
		return nil, errors.New("nat_name_not_found")
	}
	natType, ok := natConfig["type"]
	if !ok {
		log.Printf("ERROR cannot find type in the nat config")
		return nil, errors.New("missing_type")
	}
	var provider Provider
	if natType == "frp" {
		provider = &FRP{}
	} else if natType == "direct" {
		provider = &Direct{}
	} else {
		log.Printf("ERROR nat type is not supported: %v\n", natType)
		return nil, errors.New("nat_not_supported")
	}

	err := provider.InitWithConfig(natConfig)
	if err != nil {
		log.Printf("ERROR cannot init nat")
		return nil, err
	}
	gProvider = provider
	return provider, nil
}
