package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

type windowsService struct {
	stop chan struct{}
}

func (s *windowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

	for {
		select {
		case <-s.stop:
			changes <- svc.Status{State: svc.StopPending}
			return
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

func InstallService(name, desc string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", name)
	}

	config := mgr.Config{
		DisplayName:      name,
		Description:      desc,
		StartType:        mgr.StartAutomatic,
		ServiceStartName: "LocalSystem",
	}

	s, err = m.CreateService(name, exe, config)
	if err != nil {
		return err
	}
	defer s.Close()

	return nil
}

func RemoveService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", name)
	}
	defer s.Close()

	return s.Delete()
}

func RunService(name string) error {
	stop := make(chan struct{})
	service := &windowsService{stop: stop}
	return svc.Run(name, service)
} 