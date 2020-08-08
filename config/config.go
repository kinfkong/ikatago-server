package config

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/kinfkong/ikatago-server/utils"
	"github.com/spf13/viper"
)

var config *viper.Viper

// Init inits the config
func Init(configFile *string) {
	config = viper.New()

	if configFile != nil {
		config.SetConfigFile(*configFile)
	}

	config.SetDefault("world.url", utils.WorldURL)

	if configFile != nil {
		err := config.ReadInConfig()
		if err != nil {
			log.Fatal("error on parsing configuration file", err)
		}
	}
}

// GetConfig gets the configuration
func GetConfig() *viper.Viper {
	return config
}

// GetServerListenPort gets the server listen port
func GetServerListenPort() (int, error) {
	serverAddr := config.GetString("server.listen")
	items := strings.Split(serverAddr, ":")
	if len(items) < 2 {
		return 0, errors.New("invalid_server_addr")
	}
	return strconv.Atoi(items[1])
}
