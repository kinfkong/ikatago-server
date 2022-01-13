package model

type SSHLoginInfo struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	User      string `json:"user"`
	Protected bool   `json:"protected"`
}
type ServerInfoItem struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type ServerInfo struct {
	ServerVersion      string           `json:"serverVersion"`
	SupportKataWeights []ServerInfoItem `json:"supportKataWeights"`
	SupportKataNames   []ServerInfoItem `json:"supportKataNames"`
	SupportKataConfigs []ServerInfoItem `json:"supportKataConfigs"`
	DefaultKataName    string           `json:"defaultKataName"`
	DefaultKataWeight  string           `json:"defaultKataWeight"`
	DefaultKataConfig  string           `json:"defaultKataConfig"`
	GPUs               []string         `json:"gpus"`
}
