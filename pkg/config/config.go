package config

import (
<<<<<<< HEAD
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
=======
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/2gc-dev/cloudbridge-client/pkg/types"
	"github.com/spf13/viper"
)

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*types.Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/cloudbridge-client")
	viper.AddConfigPath("$HOME/.cloudbridge-client")

	// Set defaults
	setDefaults()

	// Read config file if specified
	if configPath != "" {
		viper.SetConfigFile(configPath)
	}

	// Read environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CLOUDBRIDGE")

	// Read config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var config types.Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	viper.SetDefault("relay.host", "edge.2gc.ru")
	viper.SetDefault("relay.port", 8080)
	viper.SetDefault("relay.timeout", "30s")
	viper.SetDefault("relay.tls.enabled", true)
	viper.SetDefault("relay.tls.min_version", "1.3")
	viper.SetDefault("relay.tls.verify_cert", true)
	viper.SetDefault("auth.type", "jwt")
	viper.SetDefault("auth.keycloak.enabled", false)
	viper.SetDefault("rate_limiting.enabled", true)
	viper.SetDefault("rate_limiting.max_retries", 3)
	viper.SetDefault("rate_limiting.backoff_multiplier", 2.0)
	viper.SetDefault("rate_limiting.max_backoff", "30s")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
}

// validateConfig validates the configuration
func validateConfig(c *types.Config) error {
	if c.Relay.Host == "" {
		return fmt.Errorf("relay host is required")
	}

	if c.Relay.Port <= 0 || c.Relay.Port > 65535 {
		return fmt.Errorf("invalid relay port")
	}

	if c.Relay.TLS.Enabled {
		if c.Relay.TLS.MinVersion != "1.3" {
			return fmt.Errorf("only TLS 1.3 is supported")
		}

		if c.Relay.TLS.VerifyCert {
			if c.Relay.TLS.CACert != "" {
				if _, err := os.Stat(c.Relay.TLS.CACert); os.IsNotExist(err) {
					return fmt.Errorf("CA certificate file not found: %s", c.Relay.TLS.CACert)
				}
			}
		}

		if c.Relay.TLS.ClientCert != "" && c.Relay.TLS.ClientKey == "" {
			return fmt.Errorf("client key is required when client certificate is provided")
		}

		if c.Relay.TLS.ClientKey != "" && c.Relay.TLS.ClientCert == "" {
			return fmt.Errorf("client certificate is required when client key is provided")
		}
	}

	if c.Auth.Type == "jwt" && c.Auth.Secret == "" {
		return fmt.Errorf("JWT secret is required for JWT authentication")
	}

	if c.Auth.Keycloak.Enabled {
		if c.Auth.Keycloak.ServerURL == "" {
			return fmt.Errorf("Keycloak server URL is required")
		}
		if c.Auth.Keycloak.Realm == "" {
			return fmt.Errorf("Keycloak realm is required")
		}
		if c.Auth.Keycloak.ClientID == "" {
			return fmt.Errorf("Keycloak client ID is required")
		}
	}

	if c.RateLimiting.MaxRetries < 0 {
		return fmt.Errorf("max retries cannot be negative")
	}

	if c.RateLimiting.BackoffMultiplier <= 0 {
		return fmt.Errorf("backoff multiplier must be positive")
>>>>>>> ebb63d9 (feat: implement CloudBridge Relay Client with TLS 1.3, JWT auth, tunnels, heartbeat, rate limiting, comprehensive docs and tests)
	}

	return nil
}

<<<<<<< HEAD
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
=======
// CreateTLSConfig creates a TLS configuration from the config
func CreateTLSConfig(c *types.Config) (*tls.Config, error) {
	if !c.Relay.TLS.Enabled {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_AES_128_GCM_SHA256,
		},
		InsecureSkipVerify: !c.Relay.TLS.VerifyCert,
	}

	// Set server name for SNI
	if c.Relay.TLS.ServerName != "" {
		tlsConfig.ServerName = c.Relay.TLS.ServerName
	} else {
		tlsConfig.ServerName = c.Relay.Host
	}

	// Load CA certificate if provided
	if c.Relay.TLS.CACert != "" {
		caCert, err := os.ReadFile(c.Relay.TLS.CACert)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to append CA certificate")
		}

		tlsConfig.RootCAs = caCertPool
	}

	// Load client certificate if provided
	if c.Relay.TLS.ClientCert != "" && c.Relay.TLS.ClientKey != "" {
		cert, err := tls.LoadX509KeyPair(c.Relay.TLS.ClientCert, c.Relay.TLS.ClientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return tlsConfig, nil
>>>>>>> ebb63d9 (feat: implement CloudBridge Relay Client with TLS 1.3, JWT auth, tunnels, heartbeat, rate limiting, comprehensive docs and tests)
} 