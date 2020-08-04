package model

type SSHLoginInfo struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
}
