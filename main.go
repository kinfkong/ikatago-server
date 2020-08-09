package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jessevdk/go-flags"

	"github.com/kinfkong/ikatago-server/config"
	"github.com/kinfkong/ikatago-server/katago"
	"github.com/kinfkong/ikatago-server/model"
	"github.com/kinfkong/ikatago-server/nat"
	"github.com/kinfkong/ikatago-server/platform"
	"github.com/kinfkong/ikatago-server/sshd"
	"github.com/kinfkong/ikatago-server/storage"
	"github.com/kinfkong/ikatago-server/utils"
)

const (
	ServerVersion = "1.2.1"
)

var opts struct {
	World         *string `short:"w" long:"world" description:"The world url."`
	Platform      string  `short:"p" long:"platform" description:"The platform, like aistudio, colab" required:"true"`
	PlatformToken string  `short:"t" long:"token" description:"The token of the platform, like aistudio, colab" required:"true"`
	ConfigFile    *string `short:"c" long:"config" description:"The config file of the server (not katago config file)" default:"./config/conf.yaml"`
}

func getPlatformFromWorld() (*platform.Platform, error) {
	type World struct {
		Platforms []platform.Platform `json:"platforms"`
		PublicKey string              `json:"publicKey"`
		HomePage  string              `json:"homepage"`
	}
	worldJSONString, err := utils.DoHTTPRequest("GET", config.GetConfig().GetString("world.url"), nil, nil)
	if err != nil {
		return nil, err
	}
	world := &World{}
	err = json.Unmarshal([]byte(worldJSONString), &world)
	if err != nil {
		return nil, err
	}
	for _, platform := range world.Platforms {
		if platform.Name == config.GetConfig().GetString("platform.name") {
			claims, err := validateToken(config.GetConfig().GetString("platform.token"), world.PublicKey)
			if err != nil {
				log.Fatal(err)
			}
			dataEncryptKeyPrefixV, ok := claims["dataEncryptKeyPrefix"]
			if !ok {
				log.Fatal("cannot find dataEncryptKeyPrefix")
			}
			dataEncryptKeyPrefix, ok := dataEncryptKeyPrefixV.(string)
			if !ok {
				log.Fatal("cannot find valid dataEncryptKeyPrefix")
			}
			err = platform.Decrypt(dataEncryptKeyPrefix)
			if err != nil {
				log.Fatal(err)
			}
			return &platform, nil
		}
	}
	log.Printf("ERROR platform not found in the world. platform: %s", config.GetConfig().GetString("platform.name"))
	return nil, errors.New("platform_not_found")
}

func validateToken(tokenString string, publicKey string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		pem, err := b64.StdEncoding.DecodeString(publicKey)
		if err != nil {
			return nil, err
		}
		return jwt.ParseRSAPublicKeyFromPEM(pem)
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		log.Printf("ERROR: invalid token:" + tokenString)
		return nil, errors.New("invalid_token")
	}
	if aud, _ := claims["aud"]; aud != config.GetConfig().GetString("platform.name") {
		return nil, errors.New("invalid_token")
	}
	var expV int64 = 0
	switch exp := claims["exp"].(type) {
	case float64:
		expV = int64(exp)
	case json.Number:
		expV, _ = exp.Int64()
	}
	log.Printf("Token will expires at: %v\n", time.Unix(expV, 0))
	return claims, err
}

func parseArgs() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal("Cannot parse args", err)
	}
	config.Init(opts.ConfigFile)
	// overrides the config with args
	if opts.World != nil {
		config.GetConfig().Set("world.url", *opts.World)
	}

	config.GetConfig().Set("platform.name", opts.Platform)
	config.GetConfig().Set("platform.token", opts.PlatformToken)
	log.Printf("DEBUG the world is: %s\n", config.GetConfig().GetString("world.url"))
	log.Printf("DEBUG Platform: [%s]\n", config.GetConfig().GetString("platform.name"))

}

func main() {
	fmt.Printf("Server Version: %s\n", ServerVersion)
	parseArgs()
	platform, err := getPlatformFromWorld()
	if err != nil {
		log.Fatal(err)
	}

	// check the katago manager
	katagoManager := katago.GetManager()
	if katagoManager == nil {
		log.Fatal("katago config seems wrong...")
	}

	// start sshd
	err = sshd.RunAsync()
	if err != nil {
		log.Fatal(err)
	}

	natProvider, err := nat.GetNatProvider()
	if err != nil {
		log.Fatal(err)
	}

	err = natProvider.RunAsync()
	if err != nil {
		log.Fatal(err)
	}
	sshInfo, err := natProvider.GetInfo()
	if err != nil {
		log.Fatal(err)
	}
	// upload the users
	oss := storage.Oss{
		BucketEndpoint:  platform.Oss.BucketEndpoint,
		BucketName:      platform.Oss.Bucket,
		AccessKeyId:     platform.Oss.Data.User.AccessKey,
		AccessKeySecret: platform.Oss.Data.User.AccessSecret,
	}
	oss.Init()
	for _, sshUser := range sshd.Users {
		err := oss.SaveUserSSHInfo(model.SSHLoginInfo{
			Host: sshInfo.Host,
			Port: sshInfo.Port,
			User: sshUser.Username,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("\n\n")
	fmt.Printf("SSH HOST: %s\n", sshInfo.Host)
	fmt.Printf("SSH PORT: %d\n\n", sshInfo.Port)
	fmt.Printf("\n")

	fmt.Printf("Congratulations! Now ikatago-server is running successfully, waiting for your requests ...\n\n")

	for {
		// wait for the services
		time.Sleep(1000 * time.Millisecond)
	}
}
