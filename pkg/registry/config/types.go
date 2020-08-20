package config

// Config is the config for ormb.
type Config struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Domain   string `json:"domain,omitempty"`

	KubeConfig string `json:"kube_config,omitempty"`
}

// New create a new Config.
func New() *Config {
	return &Config{}
}
