package relay

import (
	"fmt"
	"sync"
	"testing"
	"github.com/2gc-dev/cloudbridge-client/pkg/tunnel"
	"github.com/2gc-dev/cloudbridge-client/pkg/types"
)

type mockClient struct{}

func (m *mockClient) IsConnected() bool { return true }
func (m *mockClient) SendHeartbeat() error { return nil }
func (m *mockClient) GetConfig() *types.Config { return nil }
func (m *mockClient) GetClientID() string { return "mock-client" }
func (m *mockClient) GetTenantID() string { return "mock-tenant" }

func TestTunnelCreation(t *testing.T) {
	mgr := tunnel.NewManager(&mockClient{})
	err := mgr.RegisterTunnel("test-tunnel-1", 5001, "test-server", 3389)
	if err != nil {
		t.Fatalf("Failed to create tunnel: %v", err)
	}
	tun, exists := mgr.GetTunnel("test-tunnel-1")
	if !exists {
		t.Fatal("Tunnel not found after creation")
	}
	if tun.ID != "test-tunnel-1" {
		t.Errorf("Expected tunnel ID 'test-tunnel-1', got '%s'", tun.ID)
	}
	if tun.LocalPort != 5001 {
		t.Errorf("Expected local port 5001, got %d", tun.LocalPort)
	}
	if tun.RemoteHost != "test-server" {
		t.Errorf("Expected remote host 'test-server', got '%s'", tun.RemoteHost)
	}
	if tun.RemotePort != 3389 {
		t.Errorf("Expected remote port 3389, got %d", tun.RemotePort)
	}
}

func TestTunnelManagement(t *testing.T) {
	mgr := tunnel.NewManager(&mockClient{})
	tunnels := []struct {
		id         string
		localPort  int
		remoteHost string
		remotePort int
	}{
		{"test-tunnel-1", 5002, "test-server-1", 3389},
		{"test-tunnel-2", 5003, "test-server-2", 3389},
	}
	for _, info := range tunnels {
		err := mgr.RegisterTunnel(info.id, info.localPort, info.remoteHost, info.remotePort)
		if err != nil {
			t.Fatalf("Failed to create tunnel %s: %v", info.id, err)
		}
	}
	all := mgr.ListTunnels()
	if len(all) != 2 {
		t.Errorf("Expected 2 tunnels, got %d", len(all))
	}
	tun, exists := mgr.GetTunnel("test-tunnel-1")
	if !exists || tun.ID != "test-tunnel-1" {
		t.Error("Tunnel test-tunnel-1 not found")
	}
	err := mgr.UnregisterTunnel("test-tunnel-1")
	if err != nil {
		t.Errorf("Failed to close tunnel: %v", err)
	}
	all = mgr.ListTunnels()
	if len(all) != 1 {
		t.Errorf("Expected 1 tunnel after closing, got %d", len(all))
	}
	_, exists = mgr.GetTunnel("test-tunnel-1")
	if exists {
		t.Error("Closed tunnel still exists")
	}
}

func TestTunnelConcurrency(t *testing.T) {
	mgr := tunnel.NewManager(&mockClient{})
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := mgr.RegisterTunnel(fmt.Sprintf("test-tunnel-%d", id), 6000+id, "test-server", 3389)
			if err != nil {
				t.Errorf("Failed to create tunnel: %v", err)
			}
		}(i)
	}
	wg.Wait()
	all := mgr.ListTunnels()
	if len(all) != 10 {
		t.Errorf("Expected 10 tunnels, got %d", len(all))
	}
} 