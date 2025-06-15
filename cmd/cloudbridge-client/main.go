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
	"sync/atomic"

	"github.com/2gc-dev/cloudbridge-client/pkg/config"
	"github.com/2gc-dev/cloudbridge-client/pkg/relay"
)

var version = "1.0.0"

const (
	maxRetries      = 5
	initialDelaySec = 1
	maxDelaySec     = 30
)

var (
	connectionTimeSum int64
	connectionCount  int64
	errorCount        int64
	activeTunnels    int64
)

func main() {
	configPath := flag.String("config", "", "Path to config file")
	logFilePath := flag.String("logfile", "/var/log/cloudbridge-client/client.log", "Path to log file")
	flag.Parse()

	// Логирование в файл и консоль
	logFile, err := os.OpenFile(*logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

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
				atomic.AddInt64(&errorCount, 1)
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
				atomic.AddInt64(&errorCount, 1)
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

			atomic.AddInt64(&connectionCount, 1)
			atomic.AddInt64(&connectionTimeSum, int64(time.Since(start).Milliseconds()))

			err := client.EventLoop(func(tunnelInfo map[string]interface{}) {
				log.Printf("[EVENT] Tunnel registration requested: %v", tunnelInfo)
				atomic.AddInt64(&activeTunnels, 1)
				// Здесь можно реализовать динамическое создание туннеля по tunnelInfo
			})
			if err != nil {
				log.Printf("Event loop error: %v", err)
				atomic.AddInt64(&errorCount, 1)
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
			// Если EventLoop завершился без ошибки — выход
			break
		}
	}()

	// Периодический вывод метрик
	go func() {
		for {
			time.Sleep(60 * time.Second)
			cc := atomic.LoadInt64(&connectionCount)
			ct := atomic.LoadInt64(&connectionTimeSum)
			ec := atomic.LoadInt64(&errorCount)
			at := atomic.LoadInt64(&activeTunnels)
			avgConnTime := float64(ct) / float64(cc+1)
			log.Printf("[METRICS] connections: %d, avg_connection_time_ms: %.2f, errors: %d, active_tunnels: %d", cc, avgConnTime, ec, at)
		}
	}()

	<-sigChan
	log.Println("Shutting down...")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
} 