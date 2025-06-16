package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	TLS struct {
		Enabled  bool   `yaml:"enabled"`
		CertFile string `yaml:"cert_file"`
		KeyFile  string `yaml:"key_file"`
		CAFile   string `yaml:"ca_file"`
	} `yaml:"tls"`

	Server struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		JWTToken string `yaml:"jwt_token"`
	} `yaml:"server"`

	Tunnel struct {
		LocalPort      int `yaml:"local_port"`
		ReconnectDelay int `yaml:"reconnect_delay"`
		MaxRetries     int `yaml:"max_retries"`
	} `yaml:"tunnel"`

	Logging struct {
		Level      string `yaml:"level"`
		File       string `yaml:"file"`
		MaxSize    int    `yaml:"max_size"`
		MaxBackups int    `yaml:"max_backups"`
		MaxAge     int    `yaml:"max_age"`
		Compress   bool   `yaml:"compress"`
	} `yaml:"logging"`
}

// Save сохраняет конфигурацию в файл
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("error marshaling config: %v", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}

func LoadConfig(configPath string) (*Config, error) {
	// If no config path is provided, try default locations
	if configPath == "" {
		configPath = os.Getenv("CONFIG_FILE")
		if configPath == "" {
			configPath = "/etc/cloudbridge-client/config.yaml"
		}
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Parse YAML
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	// Проверяем токен в переменной окружения
	if envToken := os.Getenv("CLOUDBRIDGE_JWT_TOKEN"); envToken != "" {
		config.Server.JWTToken = envToken
	}

	// Validate config
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %v", err)
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.Server.Host == "" {
		return fmt.Errorf("server host is required")
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port")
	}
	if c.Server.JWTToken == "" {
		return fmt.Errorf("JWT token is required (set in config or CLOUDBRIDGE_JWT_TOKEN environment variable)")
	}
	if c.Tunnel.LocalPort <= 0 || c.Tunnel.LocalPort > 65535 {
		return fmt.Errorf("invalid local port")
	}
	if c.Tunnel.ReconnectDelay < 0 {
		return fmt.Errorf("reconnect delay must be positive")
	}
	if c.Tunnel.MaxRetries < 0 {
		return fmt.Errorf("max retries must be positive")
	}

	if c.TLS.Enabled {
		if c.TLS.CertFile != "" && !fileExists(c.TLS.CertFile) {
			return fmt.Errorf("TLS cert file not found: %s", c.TLS.CertFile)
		}
		if c.TLS.KeyFile != "" && !fileExists(c.TLS.KeyFile) {
			return fmt.Errorf("TLS key file not found: %s", c.TLS.KeyFile)
		}
		if c.TLS.CAFile != "" && !fileExists(c.TLS.CAFile) {
			return fmt.Errorf("TLS CA file not found: %s", c.TLS.CAFile)
		}
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
} 