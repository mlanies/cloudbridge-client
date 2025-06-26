package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/2gc-dev/cloudbridge-client/pkg/config"
	"github.com/2gc-dev/cloudbridge-client/pkg/errors"
	"github.com/2gc-dev/cloudbridge-client/pkg/relay"
	"github.com/spf13/cobra"
)

var (
	configFile string
	token      string
	tunnelID   string
	localPort  int
	remoteHost string
	remotePort int
	verbose    bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cloudbridge-client",
		Short: "CloudBridge Relay Client",
		Long:  "A cross-platform client for CloudBridge Relay with TLS 1.3 support and JWT authentication",
		RunE:  run,
	}

	// Add flags
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "Configuration file path")
	rootCmd.Flags().StringVarP(&token, "token", "t", "", "JWT token for authentication")
	rootCmd.Flags().StringVarP(&tunnelID, "tunnel-id", "i", "tunnel_001", "Tunnel ID")
	rootCmd.Flags().IntVarP(&localPort, "local-port", "l", 3389, "Local port to bind")
	rootCmd.Flags().StringVarP(&remoteHost, "remote-host", "r", "192.168.1.100", "Remote host")
	rootCmd.Flags().IntVarP(&remotePort, "remote-port", "p", 3389, "Remote port")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")

	// Mark required flags
	rootCmd.MarkFlagRequired("token")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Log platform information
	log.Printf("Running on %s/%s", runtime.GOOS, runtime.GOARCH)

	// Load configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override config with command line flags if provided
	if token != "" {
		cfg.Auth.Secret = token // For JWT auth, secret is the token
	}

	// Create client
	client, err := relay.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	// Set up signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	if runtime.GOOS == "windows" {
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	} else {
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	}

	// Start connection with retry logic
	if err := connectWithRetry(client); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	log.Printf("Successfully connected to relay server %s:%d", cfg.Relay.Host, cfg.Relay.Port)

	// Authenticate
	if err := authenticateWithRetry(client, token); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	log.Printf("Successfully authenticated with client ID: %s", client.GetClientID())

	// Create tunnel
	if err := createTunnelWithRetry(client, tunnelID, localPort, remoteHost, remotePort); err != nil {
		return fmt.Errorf("failed to create tunnel: %w", err)
	}

	log.Printf("Successfully created tunnel %s: localhost:%d -> %s:%d", 
		tunnelID, localPort, remoteHost, remotePort)

	// Start heartbeat
	if err := client.StartHeartbeat(); err != nil {
		return fmt.Errorf("failed to start heartbeat: %w", err)
	}

	log.Printf("Heartbeat started")

	// Wait for shutdown signal
	select {
	case <-sigChan:
		log.Println("Received shutdown signal, closing...")
	case <-ctx.Done():
		log.Println("Context cancelled, closing...")
	}

	return nil
}

// connectWithRetry connects to the relay server with retry logic
func connectWithRetry(client *relay.Client) error {
	retryStrategy := client.GetRetryStrategy()
	
	for {
		err := client.Connect()
		if err == nil {
			return nil
		}

		relayErr, _ := errors.HandleError(err)
		if relayErr == nil || !retryStrategy.ShouldRetry(err) {
			return err
		}

		delay := retryStrategy.GetNextDelay(err)
		log.Printf("Connection failed: %v, retrying in %v...", err, delay)
		time.Sleep(delay)
	}
}

// authenticateWithRetry authenticates with retry logic
func authenticateWithRetry(client *relay.Client, token string) error {
	retryStrategy := client.GetRetryStrategy()
	
	for {
		err := client.Authenticate(token)
		if err == nil {
			return nil
		}

		relayErr, _ := errors.HandleError(err)
		if relayErr == nil || !retryStrategy.ShouldRetry(err) {
			return err
		}

		delay := retryStrategy.GetNextDelay(err)
		log.Printf("Authentication failed: %v, retrying in %v...", err, delay)
		time.Sleep(delay)
	}
}

// createTunnelWithRetry creates a tunnel with retry logic
func createTunnelWithRetry(client *relay.Client, tunnelID string, localPort int, remoteHost string, remotePort int) error {
	retryStrategy := client.GetRetryStrategy()
	
	for {
		err := client.CreateTunnel(tunnelID, localPort, remoteHost, remotePort)
		if err == nil {
			return nil
		}

		relayErr, _ := errors.HandleError(err)
		if relayErr == nil || !retryStrategy.ShouldRetry(err) {
			return err
		}

		delay := retryStrategy.GetNextDelay(err)
		log.Printf("Tunnel creation failed: %v, retrying in %v...", err, delay)
		time.Sleep(delay)
	}
} 