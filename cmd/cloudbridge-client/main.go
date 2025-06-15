package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/yourusername/cloudbridge-client/pkg/relay"
)

func main() {
	// Parse command line flags
	config := &relay.Config{}
	flag.BoolVar(&config.UseTLS, "tls", true, "Use TLS for connection")
	flag.StringVar(&config.TLSCertFile, "cert", "", "TLS certificate file")
	flag.StringVar(&config.TLSKeyFile, "key", "", "TLS key file")
	flag.StringVar(&config.TLSCAFile, "ca", "", "TLS CA certificate file")
	flag.StringVar(&config.ServerHost, "host", "edge.2gc.ru", "Relay server host")
	flag.IntVar(&config.ServerPort, "port", 8080, "Relay server port")
	flag.StringVar(&config.JWTToken, "token", "", "JWT token for authentication")
	flag.IntVar(&config.LocalPort, "local-port", 3389, "Local port to tunnel")
	flag.IntVar(&config.ReconnectDelay, "reconnect-delay", 5, "Reconnection delay in seconds")
	flag.IntVar(&config.MaxRetries, "max-retries", 3, "Maximum number of reconnection attempts")
	flag.Parse()

	// Log platform information
	log.Printf("Running on %s/%s", runtime.GOOS, runtime.GOARCH)

	// Validate configuration
	if err := config.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Create TLS config if needed
	var tlsConfig *tls.Config
	var err error
	if config.UseTLS {
		tlsConfig, err = relay.NewTLSConfig(config.TLSCertFile, config.TLSKeyFile, config.TLSCAFile)
		if err != nil {
			log.Fatalf("Failed to create TLS config: %v", err)
		}
	}

	// Create client
	client := relay.NewClient(config.UseTLS, tlsConfig)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	if runtime.GOOS == "windows" {
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	} else {
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	}

	// Connect to relay server
	if err := client.Connect(config.ServerHost, config.ServerPort); err != nil {
		log.Fatalf("Failed to connect to relay server: %v", err)
	}
	defer client.Close()

	// Register with JWT token
	if err := client.Register(config.JWTToken); err != nil {
		log.Fatalf("Failed to register with relay server: %v", err)
	}

	// Create tunnel
	if err := client.CreateTunnel(config.LocalPort); err != nil {
		log.Fatalf("Failed to create tunnel: %v", err)
	}

	log.Printf("Successfully connected to relay server and created tunnel on port %d", config.LocalPort)

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down...")
} 