package nat

// Info the ssh info
type Info struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// Provider the interface of the provider
type Provider interface {
	RunAsync() error
	GetInfo() (Info, error)
}
