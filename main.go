package main

import (
	"fmt"
	"log"
	"os"

	"redis-internal/config"
	"redis-internal/core"
	"redis-internal/server"
)

func main() {
	fmt.Println("Starting the Redis Internal server...")

	// Load configuration from file and command line
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := appConfig.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Print configuration
	appConfig.Print()

	// Initialize the core store with configuration
	storeConfig := core.StoreConfig{
		KeysLimit:        appConfig.KeysLimit,
		EvictionStrategy: appConfig.EvictionStrategy,
	}
	core.InitStore(storeConfig)

	// Convert to server.Config type
	serverConfig := server.Config{
		Host:                appConfig.Host,
		Port:                appConfig.Port,
		KeysLimit:           appConfig.KeysLimit,
		EvictionStrategy:    appConfig.EvictionStrategy,
		AutoDeleteFrequency: appConfig.AutoDeleteFrequency,
		MaxClients:          appConfig.MaxClients,
		LogLevel:            appConfig.LogLevel,
	}

	// Start the async TCP server
	if err := server.RunAsyncTCPServer(serverConfig); err != nil {
		log.Fatalf("Server error: %v", err)
		os.Exit(1)
	}
}
