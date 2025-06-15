package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
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

	sigChan := make(chan os.Signal, 1)
	if runtime.GOOS == "windows" {
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	} else {
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	}

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
} 