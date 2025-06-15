package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"crypto/tls"

	"github.com/2gc-dev/cloudbridge-client/pkg/config"
	"github.com/2gc-dev/cloudbridge-client/pkg/relay"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Log platform information
	log.Printf("Running on %s/%s", runtime.GOOS, runtime.GOARCH)

	// Create TLS config if needed
	var tlsConfig *tls.Config
	if cfg.TLS.Enabled {
		tlsConfig, err = relay.NewTLSConfig(cfg.TLS.CertFile, cfg.TLS.KeyFile, cfg.TLS.CAFile)
		if err != nil {
			log.Fatalf("Failed to create TLS config: %v", err)
		}
	}

	// Create client
	client := relay.NewClient(cfg.TLS.Enabled, tlsConfig)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	if runtime.GOOS == "windows" {
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	} else {
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	}

	// Connect to relay server
	if err := client.Connect(cfg.Server.Host, cfg.Server.Port); err != nil {
		log.Fatalf("Failed to connect to relay server: %v", err)
	}
	defer client.Close()

	// Register with JWT token
	if err := client.Register(cfg.Server.JWTToken); err != nil {
		log.Fatalf("Failed to register with relay server: %v", err)
	}

	// Create tunnel
	if err := client.CreateTunnel(cfg.Tunnel.LocalPort); err != nil {
		log.Fatalf("Failed to create tunnel: %v", err)
	}

	log.Printf("Successfully connected to relay server and created tunnel on port %d", cfg.Tunnel.LocalPort)

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down...")
} 