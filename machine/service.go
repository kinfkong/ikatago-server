package machine

import (
	b64 "encoding/base64"

	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Service the machine service
type Service struct {
	MachineID string `json:"machineId"`
	PubKey    string `json:"pubKey"`
	NatHost   string `json:"natHost"`
	NatPort   string `json:"natPort"`
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

func (service *Service) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		pem, err := b64.StdEncoding.DecodeString(service.PubKey)
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
	if sub, ok := claims["sub"].(string); !ok || len(sub) == 0 {
		return nil, errors.New("invalid_token")
	}
	var expV int64 = 0
	switch exp := claims["exp"].(type) {
	case float64:
		expV = int64(exp)
	case json.Number:
		expV, _ = exp.Int64()
	}
	log.Printf("Machine Token will expires at: %v\n", time.Unix(expV, 0))
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
