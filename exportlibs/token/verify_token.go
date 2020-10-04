package main

// #cgo LDFLAGS:
import "C"
import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kinfkong/ikatago-server/platform"
	"github.com/kinfkong/ikatago-server/utils"
)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

//export VerifySpecificPlatformToken
func VerifySpecificPlatformToken(platformNameHash uint32) int {
	if hash(os.Getenv("IKATAGO_PLATFORM_NAME")) != platformNameHash {
		return -1
	}
	return VerifyPlatformToken()
}

//export VerifyPlatformToken
func VerifyPlatformToken() int {
	platformName := os.Getenv("IKATAGO_PLATFORM_NAME")
	platformToken := os.Getenv("IKATAGO_PLATFORM_TOKEN")
	type World struct {
		Platforms []platform.Platform `json:"platforms"`
		PublicKey string              `json:"publicKey"`
		HomePage  string              `json:"homepage"`
	}
	worldJSONString, err := utils.DoHTTPRequest("GET", utils.WorldURL, nil, nil)
	if err != nil {
		return -1
	}
	world := &World{}
	err = json.Unmarshal([]byte(worldJSONString), &world)
	if err != nil {
		return -1
	}
	for _, platform := range world.Platforms {
		if platform.Name == platformName {
			claims, err := validateToken(platformName, platformToken, world.PublicKey)
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
			return 0
		}
	}
	log.Printf("ERROR platform not found in the world. platform: %s", platformName)
	return -1
}

func validateToken(platformName string, tokenString string, publicKey string) (jwt.MapClaims, error) {
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
	if aud, _ := claims["aud"]; aud != platformName {
		return nil, errors.New("invalid_token")
	}
	var expV int64 = 0
	switch exp := claims["exp"].(type) {
	case float64:
		expV = int64(exp)
	case json.Number:
		expV, _ = exp.Int64()
	}
	// log.Printf("Token will expires at: %v\n", time.Unix(expV, 0))
	// use a thread to verify the token
	expireDate := time.Unix(expV, 0)
	go func() {
		for {
			now := time.Now()
			if now.After(expireDate) {
				log.Fatal("Token expired")
			}
			time.Sleep(5 * time.Second)
		}
	}()
	return claims, err
}

func main() {
	fmt.Printf("%v\n", hash("aistudio-8v"))
}
