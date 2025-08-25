package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

// AppConfig represents the application configuration
type AppConfig struct {
	Host                string `json:"host"`
	Port                int    `json:"port"`
	KeysLimit           int    `json:"keysLimit"`
	EvictionStrategy    string `json:"evictionStrategy"`
	AutoDeleteFrequency string `json:"autoDeleteFrequency"`
	MaxClients          int    `json:"maxClients"`
	LogLevel            string `json:"logLevel"`
}

// DefaultConfig returns default configuration values
func DefaultConfig() *AppConfig {
	return &AppConfig{
		Host:                "0.0.0.0",
		Port:                7379,
		KeysLimit:           1000,
		EvictionStrategy:    "simple-first",
		AutoDeleteFrequency: "1s",
		MaxClients:          20000,
		LogLevel:            "info",
	}
}

// LoadConfig loads configuration from file with command line flag overrides
func LoadConfig() (*AppConfig, error) {
	var configFile string
	flag.StringVar(&configFile, "config", "config.json", "path to configuration file")

	// Create default config
	config := DefaultConfig()

	// Define command line flags that can override config file
	var (
		host             = flag.String("host", "", "host for the redis server")
		port             = flag.Int("port", 0, "port for the redis server")
		keysLimit        = flag.Int("keys-limit", 0, "maximum key limit")
		evictionStrategy = flag.String("eviction", "", "eviction strategy (simple-first)")
		maxClients       = flag.Int("max-clients", 0, "maximum number of clients")
		logLevel         = flag.String("log-level", "", "log level (info, debug, warn, error)")
	)

	flag.Parse()

	// Load from config file if it exists
	if _, err := os.Stat(configFile); err == nil {
		file, err := os.Open(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %v", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %v", err)
		}

		fmt.Printf("Loaded configuration from %s\n", configFile)
	} else {
		fmt.Println("Config file not found, using defaults and command line arguments")
	}

	// Override with command line flags if provided
	if *host != "" {
		config.Host = *host
	}
	if *port != 0 {
		config.Port = *port
	}
	if *keysLimit != 0 {
		config.KeysLimit = *keysLimit
	}
	if *evictionStrategy != "" {
		config.EvictionStrategy = *evictionStrategy
	}
	if *maxClients != 0 {
		config.MaxClients = *maxClients
	}
	if *logLevel != "" {
		config.LogLevel = *logLevel
	}

	return config, nil
}

// GetAutoDeleteDuration parses the AutoDeleteFrequency and returns a time.Duration
func (c *AppConfig) GetAutoDeleteDuration() (time.Duration, error) {
	return time.ParseDuration(c.AutoDeleteFrequency)
}

// Validate checks if the configuration values are valid
func (c *AppConfig) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", c.Port)
	}

	if c.KeysLimit < 1 {
		return fmt.Errorf("keys limit must be greater than 0: %d", c.KeysLimit)
	}

	if c.MaxClients < 1 {
		return fmt.Errorf("max clients must be greater than 0: %d", c.MaxClients)
	}

	// Validate auto-delete frequency
	if _, err := c.GetAutoDeleteDuration(); err != nil {
		return fmt.Errorf("invalid auto delete frequency: %v", err)
	}

	// Validate eviction strategy
	validStrategies := []string{"simple-first", "lru", "random"}
	valid := false
	for _, strategy := range validStrategies {
		if c.EvictionStrategy == strategy {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid eviction strategy: %s", c.EvictionStrategy)
	}

	return nil
}

// Print displays the current configuration
func (c *AppConfig) Print() {
	fmt.Println("=== Redis Internal Configuration ===")
	fmt.Printf("Host: %s\n", c.Host)
	fmt.Printf("Port: %d\n", c.Port)
	fmt.Printf("Keys Limit: %d\n", c.KeysLimit)
	fmt.Printf("Eviction Strategy: %s\n", c.EvictionStrategy)
	fmt.Printf("Auto Delete Frequency: %s\n", c.AutoDeleteFrequency)
	fmt.Printf("Max Clients: %d\n", c.MaxClients)
	fmt.Printf("Log Level: %s\n", c.LogLevel)
	fmt.Println("===================================")
}
