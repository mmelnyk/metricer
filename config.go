package metricer

// Config represents configuration structure
type Config struct {
	// AllowExternal
	AllowExternal bool `json:"external,omit" yaml:"external"`

	// EnableDebug
	EnableDebug bool `json:"debug,omit" yaml:"debug"`

	// Port
	Port uint `jsonL:"port, omit" yaml:"port"`
}

// Validate checks config structure
func (config *Config) Validate() error {
	if config == nil {
		return errNilConfig
	}

	if config.Port == 0 {
		config.Port = defaultPort
	}
	return nil
}
