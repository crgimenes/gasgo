package config

import "github.com/crgimenes/goconfig"

// Config holds the configuration.
type Config struct {
	MetadataFile string `json:"metadata_file" cfg:"metadata_file"`
}

// Parse condig parameters
func Parse() (*Config, error) {
	cfg := &Config{}

	err := goconfig.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
