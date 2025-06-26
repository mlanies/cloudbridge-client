package main

import (
<<<<<<< HEAD
	"flag"
	"io"
=======
	"context"
	"flag"
	"fmt"
>>>>>>> ebb63d9 (feat: implement CloudBridge Relay Client with TLS 1.3, JWT auth, tunnels, heartbeat, rate limiting, comprehensive docs and tests)
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
<<<<<<< HEAD
	"crypto/tls"
	"time"
	"net/http"

	"github.com/2gc-dev/cloudbridge-client/pkg/config"
	"github.com/2gc-dev/cloudbridge-client/pkg/relay"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var version = "1.0.0"

const (
	maxRetries      = 5
	initialDelaySec = 1
	maxDelaySec     = 30
)

func main() {
	// Если есть аргументы командной строки, обрабатываем их как команды
	if len(os.Args) > 1 {
		if err := parseCommand(); err != nil {
			log.Fatalf("Command error: %v", err)
		}
		return
	}

	// Оригинальные флаги
	configPath := flag.String("config", "", "Path to config file")
	logFilePath := flag.String("logfile", "/var/log/cloudbridge-client/client.log", "Path to log file")
	metricsAddr := flag.String("metrics-addr", ":9090", "Address to serve metrics on")
	flag.Parse()

	// Логирование в файл и консоль
	logFile, err := os.OpenFile(*logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	// Запуск метрик и health check
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.Handle("/health", http.HandlerFunc(relay.HealthCheckHandler))
		if err := http.ListenAndServe(*metricsAddr, nil); err != nil {
			log.Printf("Failed to start metrics server: %v", err)
		}
	}()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Running on %s/%s", runtime.GOOS, runtime.GOARCH)

	var tlsConfig *tls.Config
	if cfg.TLS.Enabled {
		tlsConfig, err = relay.NewTLSConfig(cfg.TLS.CertFile, cfg.TLS.KeyFile, cfg.TLS.CAFile)
		if err != nil {
			log.Fatalf("Failed to create TLS config: %v", err)
		}
	}

=======
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

>>>>>>> ebb63d9 (feat: implement CloudBridge Relay Client with TLS 1.3, JWT auth, tunnels, heartbeat, rate limiting, comprehensive docs and tests)
	sigChan := make(chan os.Signal, 1)
	if runtime.GOOS == "windows" {
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	} else {
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	}

<<<<<<< HEAD
	go func() {
		retries := 0
		delay := initialDelaySec
		for {
			start := time.Now()
			client := relay.NewClient(cfg.TLS.Enabled, tlsConfig)
			if err := client.Connect(cfg.Server.Host, cfg.Server.Port); err != nil {
				log.Printf("Failed to connect to relay server: %v", err)
				relay.RecordError("connection_failed")
				relay.UpdateHealthStatus("degraded")
				retries++
				if retries > maxRetries {
					log.Fatalf("Max reconnect attempts reached. Exiting.")
				}
				log.Printf("Retrying in %d seconds...", delay)
				time.Sleep(time.Duration(delay) * time.Second)
				delay = min(delay*2, maxDelaySec)
				continue
			}
			retries = 0
			delay = initialDelaySec
			defer client.Close()

			if err := client.Handshake(cfg.Server.JWTToken, version); err != nil {
				log.Printf("Handshake failed: %v", err)
				relay.RecordError("handshake_failed")
				relay.UpdateHealthStatus("degraded")
				client.Close()
				retries++
				if retries > maxRetries {
					log.Fatalf("Max reconnect attempts reached. Exiting.")
				}
				log.Printf("Retrying in %d seconds...", delay)
				time.Sleep(time.Duration(delay) * time.Second)
				delay = min(delay*2, maxDelaySec)
				continue
			}

			relay.RecordConnection(time.Since(start).Seconds())
			relay.UpdateHealthStatus("ok")

			err := client.EventLoop(func(tunnelInfo map[string]interface{}) {
				log.Printf("[EVENT] Tunnel registration requested: %v", tunnelInfo)
				_, err := client.CreateTunnel(tunnelInfo)
				if err != nil {
					log.Printf("Failed to create tunnel: %v", err)
					relay.RecordError("tunnel_creation_failed")
					return
				}
				relay.SetActiveTunnels(len(client.ListTunnels()))
			})
			if err != nil {
				log.Printf("Event loop error: %v", err)
				relay.RecordError("event_loop_failed")
				relay.UpdateHealthStatus("degraded")
				client.Close()
				retries++
				if retries > maxRetries {
					log.Fatalf("Max reconnect attempts reached. Exiting.")
				}
				log.Printf("Retrying in %d seconds...", delay)
				time.Sleep(time.Duration(delay) * time.Second)
				delay = min(delay*2, maxDelaySec)
				continue
			}
			break
		}
	}()

	<-sigChan
	log.Println("Shutting down...")
	relay.UpdateHealthStatus("shutting_down")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
=======
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
>>>>>>> ebb63d9 (feat: implement CloudBridge Relay Client with TLS 1.3, JWT auth, tunnels, heartbeat, rate limiting, comprehensive docs and tests)
} 