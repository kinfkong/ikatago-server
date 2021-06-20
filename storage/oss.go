package storage

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	aliynOss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/kinfkong/ikatago-server/model"
)

type Oss struct {
	BucketEndpoint string `json:"bucketEndpoint"`
	BucketName     string `json:"bucketName"`

	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`

	client *aliynOss.Client
}

// Init inits it
func (oss *Oss) Init() error {
	// ensure to use https endpoint
	client, err := aliynOss.New("https://"+oss.BucketEndpoint, oss.AccessKeyId, oss.AccessKeySecret)
	if err != nil {
		return err
	}
	oss.client = client
	return nil
}

func (oss *Oss) SaveUserSSHInfo(userSSHInfo model.SSHLoginInfo) error {
	fileKey := "users/" + userSSHInfo.User + ".ssh.json"
	bucket, err := oss.client.Bucket(oss.BucketName)
	if err != nil {
		return err
	}
	readReader, err := bucket.GetObject(fileKey)
	found := false
	if err != nil {
		if !strings.Contains(err.Error(), "ErrorCode=NoSuchKey") {
			log.Printf("ERROR failed to save user to oss: %s, %+v", userSSHInfo.User, err)
			return err
		}
		// not found
	} else {
		found = true
	}
	if readReader != nil {
		defer readReader.Close()
	}
	if found {
		existingBytes, err := ioutil.ReadAll(readReader)
		if err != nil {
			return err
		}
		existingData := model.SSHLoginInfo{}
		err = json.Unmarshal(existingBytes, &existingData)
		if err != nil {
			return err
		}
		if existingData.Protected {
			log.Printf("")
			log.Printf("==== WARNING ====")
			log.Printf("WARN: user [%s] is PROTECTED in other platform!!! You cannot use this username, please change to another username.", userSSHInfo.User)
			log.Printf("==== WARNING ====")
			return nil
		}
	}
	userSSHInfo.Protected = false
	jsonV, err := json.Marshal(userSSHInfo)
	if err != nil {
		return err
	}

	reader := strings.NewReader(string(jsonV))

	err = bucket.PutObject(fileKey, reader, aliynOss.ACL(aliynOss.ACLPublicRead))
	if err != nil {
		return err
	}
	err = bucket.SetObjectACL(fileKey, aliynOss.ACLPublicRead)
	if err != nil {
		return err
	}
	log.Printf("done: %s", userSSHInfo.User)
	return nil
}
