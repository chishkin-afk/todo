package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

type loader struct {
	validate *validator.Validate
}

func New() *loader {
	return &loader{
		validate: validator.New(),
	}
}

func (l *loader) Init(path string) (*Config, error) {
	bytes, err := l.parseFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config structure: %w", err)
	}

	if err := l.validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return &cfg, nil
}

func (l *loader) validateConfig(cfg *Config) error {
	if err := l.validate.Struct(cfg); err != nil {
		return err
	}

	if cfg.Server.HTTP.TLS.Enable {
		if _, err := os.Stat(cfg.Server.HTTP.TLS.ServerCertPath); os.IsNotExist(err) {
			return err
		}

		if _, err := os.Stat(cfg.Server.HTTP.TLS.ServerKeyPath); os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

// parseFile realizes mechanic of reading file and expanding env variables
func (l *loader) parseFile(path string) ([]byte, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := os.ExpandEnv(string(bytes))

	return []byte(content), nil
}
