package platform

import (
	"encoding/json"
	"errors"

	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"

	"github.com/mergermarket/go-pkcs7"
)

// RAMUser represents the oss bucket of this platform
type RAMUser struct {
	AccessKey    string `json:"accessKey"`
	AccessSecret string `json:"accessSecret"`
}

// OssData represents the oss data
type OssData struct {
	User RAMUser `json:"user"`
}

// Oss represents the the oss
type Oss struct {
	BucketEndpoint string  `json:"bucketEndpoint"`
	Bucket         string  `json:"bucket"`
	Data           OssData `json:"data"`
	EncryptedData  string  `json:"encryptedData"`
}

// Platform represents the platform
type Platform struct {
	Name string `json:"name"`
	Oss  Oss    `json:"oss"`
}

func decrypt(encrypted string, secretKey string) (string, error) {
	padded := fmt.Sprintf("%-32v", secretKey)
	key := []byte(padded)
	cipherText, _ := hex.DecodeString(encrypted)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	if len(cipherText)%aes.BlockSize != 0 {
		return "", errors.New("cipherText is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	cipherText, _ = pkcs7.Unpad(cipherText, aes.BlockSize)
	return string(cipherText), nil
}

// Decrypt the platform data
func (p *Platform) Decrypt(dataEncryptKeyPrefix string) error {
	dataEncryptKey := dataEncryptKeyPrefix + "kjdks2ikdjskfjdks"
	decrypted, err := decrypt(p.Oss.EncryptedData, dataEncryptKey)
	ossData := OssData{}
	err = json.Unmarshal([]byte(decrypted), &ossData)
	if err != nil {
		return err
	}
	p.Oss.Data = ossData
	return err
}
