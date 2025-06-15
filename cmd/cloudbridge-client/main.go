package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"crypto/tls"
	"time"
	"net/http"

	"github.com/2gc-dev/cloudbridge-client/pkg/config"
	"github.com/2gc-dev/cloudbridge-client/pkg/relay"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

var (
	version = "dev"
	serviceName = "CloudBridgeClient"
	serviceDesc = "CloudBridge Client Service"
)

const (
	maxRetries      = 5
	initialDelaySec = 1
	maxDelaySec     = 30
)

func main() {
	// Parse command line arguments
	install := flag.Bool("install", false, "Install Windows service")
	uninstall := flag.Bool("uninstall", false, "Uninstall Windows service")
	configPath := flag.String("config", "", "Path to config file")
	logFilePath := flag.String("logfile", "/var/log/cloudbridge-client/client.log", "Path to log file")
	metricsAddr := flag.String("metrics-addr", ":9090", "Address to serve metrics on")
	flag.Parse()

	// Check if running as a service
	isService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("Failed to determine if running as service: %v", err)
	}

	if isService {
		// Run as a Windows service
		err = svc.Run(serviceName, &windowsService{})
		if err != nil {
			log.Fatalf("Service failed: %v", err)
		}
		return
	}

	// Handle service installation/uninstallation
	if *install {
		err = installService()
		if err != nil {
			log.Fatalf("Failed to install service: %v", err)
		}
		fmt.Println("Service installed successfully")
		return
	}

	if *uninstall {
		err = uninstallService()
		if err != nil {
			log.Fatalf("Failed to uninstall service: %v", err)
		}
		fmt.Println("Service uninstalled successfully")
		return
	}

	// Run as a regular application
	runApplication()
}

func installService() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", serviceName)
	}

	config := mgr.Config{
		DisplayName:      serviceName,
		Description:      serviceDesc,
		StartType:        mgr.StartAutomatic,
		ServiceStartName: "LocalSystem",
	}

	s, err = m.CreateService(serviceName, exe, config)
	if err != nil {
		return err
	}
	defer s.Close()

	return nil
}

func uninstallService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("service %s is not installed", serviceName)
	}
	defer s.Close()

	return s.Delete()
}

type windowsService struct{}

func (s *windowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

	// Start the application in a goroutine
	go runApplication()

	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				changes <- svc.Status{State: svc.StopPending}
				return
			default:
				log.Printf("unexpected control request #%d", c)
			}
		}
	}
}

func runApplication() {
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