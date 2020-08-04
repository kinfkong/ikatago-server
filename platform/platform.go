package platform

// RAMUser represents the oss bucket of this platform
type RAMUser struct {
	AccessKey    string `json:"accessKey"`
	AccessSecret string `json:"accessSecret"`
}

// Oss represents the the oss
type Oss struct {
	BucketEndpoint string `json:"bucketEndpoint"`
	Data           struct {
		User RAMUser `json:"user"`
	} `json:"data"`
}

// Platform represents the platform
type Platform struct {
	Name string `json:"name"`
	Oss  Oss    `json:"oss"`
}
