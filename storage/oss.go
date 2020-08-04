package storage

import (
	"encoding/json"
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
	client, err := aliynOss.New(oss.BucketEndpoint, oss.AccessKeyId, oss.AccessKeySecret)
	if err != nil {
		return err
	}
	oss.client = client
	return nil
}

func (oss *Oss) SaveUserSSHInfo(userSSHInfo model.SSHLoginInfo) error {
	bucket, err := oss.client.Bucket(oss.BucketName)
	if err != nil {
		return err
	}

	jsonV, err := json.Marshal(userSSHInfo)
	if err != nil {
		return err
	}

	reader := strings.NewReader(string(jsonV))

	err = bucket.PutObject("users/"+userSSHInfo.User+".ssh.json", reader)
	if err != nil {
		return err
	}

	return nil
}
