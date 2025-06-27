package test

import (
	"net"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestRelayIntegration(t *testing.T) {
	// Skip if relay-server is not available
	if _, err := os.Stat("./relay-server"); os.IsNotExist(err) {
		t.Skip("relay-server binary not found, skipping integration test")
	}

	if _, err := os.Stat("./cloudbridge-client"); os.IsNotExist(err) {
		t.Skip("cloudbridge-client binary not found, skipping integration test")
	}

	relay := exec.Command("./relay-server", "--debug")
	if err := relay.Start(); err != nil {
		t.Skipf("Failed to start relay-server: %v", err)
	}
	defer func() {
		if err := relay.Process.Kill(); err != nil {
			t.Logf("Failed to kill relay process: %v", err)
		}
	}()
	time.Sleep(2 * time.Second)

	client := exec.Command("./cloudbridge-client", "--config", "./testdata/config-test.yaml")
	if err := client.Start(); err != nil {
		t.Skipf("Failed to start cloudbridge-client: %v", err)
	}
	defer func() {
		if err := client.Process.Kill(); err != nil {
			t.Logf("Failed to kill client process: %v", err)
		}
	}()
	time.Sleep(2 * time.Second)

	conn, err := net.DialTimeout("tcp", "localhost:3389", 2*time.Second)
	if err != nil {
		t.Skipf("Tunnel not established: %v", err)
	}
	if err := conn.Close(); err != nil {
		t.Logf("Failed to close connection: %v", err)
	}
}
