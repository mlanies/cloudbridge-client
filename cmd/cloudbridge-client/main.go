//go:build windows
// +build windows

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/2gc-dev/cloudbridge-client/pkg/config"
	"github.com/2gc-dev/cloudbridge-client/pkg/relay"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

var (
	version     = "dev"
	serviceName = "CloudBridgeClient"
	serviceDesc = "CloudBridge Client Service"
)

const (
	defaultConfigPath  = "C:\\Program Files\\CloudBridge Client\\config.yaml"
	defaultLogPath     = "C:\\Program Files\\CloudBridge Client\\logs\\client.log"
	defaultMetricsAddr = ":9090"
)

type windowsService struct {
	configPath  string
	logFilePath string
	metricsAddr string
}

func (s *windowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	changes <- svc.Status{State: svc.StartPending}
	
	// Логирование в файл и консоль
	logFile, err := os.OpenFile(s.logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Настройка логирования
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Загрузка конфигурации
	cfg, err := config.LoadConfig(s.configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создание и запуск relay
	client := relay.NewClient(cfg.TLS.Enabled, nil)
	if err := client.Connect(cfg.Server.Host, cfg.Server.Port); err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer client.Close()

	// Запуск метрик
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(s.metricsAddr, nil); err != nil {
			log.Printf("Failed to start metrics server: %v", err)
		}
	}()

	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				changes <- svc.Status{State: svc.StopPending}
				client.Close()
				return
			default:
				log.Printf("Unexpected control request #%d", c)
			}
		}
	}
}

func main() {
	// Parse command line arguments
	install := flag.Bool("install", false, "Install Windows service")
	uninstall := flag.Bool("uninstall", false, "Uninstall Windows service")
	configPath := flag.String("config", defaultConfigPath, "Path to config file")
	logFilePath := flag.String("logfile", defaultLogPath, "Path to log file")
	metricsAddr := flag.String("metrics-addr", defaultMetricsAddr, "Address to serve metrics on")
	flag.Parse()

	// Check if running as a service
	isService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("Failed to determine if running as a service: %v", err)
	}

	if isService {
		// Run as a Windows service
		service := &windowsService{
			configPath:  *configPath,
			logFilePath: *logFilePath,
			metricsAddr: *metricsAddr,
		}
		err = svc.Run("CloudBridgeClient", service)
		if err != nil {
			log.Fatalf("Service failed: %v", err)
		}
	} else {
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
		runApplication(*configPath, *logFilePath, *metricsAddr)
	}
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

func runApplication(configPath, logFilePath, metricsAddr string) {
	// Логирование в файл и консоль
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Настройка логирования
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Загрузка конфигурации
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создание и запуск relay
	client := relay.NewClient(cfg.TLS.Enabled, nil)
	if err := client.Connect(cfg.Server.Host, cfg.Server.Port); err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer client.Close()

	// Запуск метрик
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(metricsAddr, nil); err != nil {
			log.Printf("Failed to start metrics server: %v", err)
		}
	}()

	// Обработка сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Ожидание сигнала завершения
	<-sigChan
	log.Println("Shutting down...")
	client.Close()
} 