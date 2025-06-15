package relay

import (
	"context"
	"fmt"
	"testing"
)

func TestTunnelCreation(t *testing.T) {
	client := NewClient(false, nil)
	
	tunnelInfo := map[string]interface{}{
		"tunnel_id":    "test-tunnel-1",
		"local_port":   3389,
		"remote_host":  "test-server",
		"remote_port":  3389,
		"protocol":     "rdp",
		"options": map[string]interface{}{
			"username": "test",
			"password": "test",
			"domain":   "test",
		},
	}
	
	tunnel, err := client.CreateTunnel(tunnelInfo)
	if err != nil {
		t.Errorf("Failed to create tunnel: %v", err)
	}
	
	if tunnel.ID != "test-tunnel-1" {
		t.Errorf("Expected tunnel ID 'test-tunnel-1', got '%s'", tunnel.ID)
	}
	
	if tunnel.LocalPort != 3389 {
		t.Errorf("Expected local port 3389, got %d", tunnel.LocalPort)
	}
	
	if tunnel.RemoteHost != "test-server" {
		t.Errorf("Expected remote host 'test-server', got '%s'", tunnel.RemoteHost)
	}
	
	if tunnel.RemotePort != 3389 {
		t.Errorf("Expected remote port 3389, got %d", tunnel.RemotePort)
	}
	
	if tunnel.Protocol != "rdp" {
		t.Errorf("Expected protocol 'rdp', got '%s'", tunnel.Protocol)
	}
}

func TestTunnelManagement(t *testing.T) {
	client := NewClient(false, nil)
	
	// Создаем несколько туннелей
	tunnelInfos := []map[string]interface{}{
		{
			"tunnel_id":    "test-tunnel-1",
			"local_port":   3389,
			"remote_host":  "test-server-1",
			"remote_port":  3389,
			"protocol":     "rdp",
			"options": map[string]interface{}{
				"username": "test1",
				"password": "test1",
				"domain":   "test1",
			},
		},
		{
			"tunnel_id":    "test-tunnel-2",
			"local_port":   3390,
			"remote_host":  "test-server-2",
			"remote_port":  3389,
			"protocol":     "rdp",
			"options": map[string]interface{}{
				"username": "test2",
				"password": "test2",
				"domain":   "test2",
			},
		},
	}
	
	// Создаем туннели
	for _, info := range tunnelInfos {
		_, err := client.CreateTunnel(info)
		if err != nil {
			t.Errorf("Failed to create tunnel %s: %v", info["tunnel_id"], err)
		}
	}
	
	// Проверяем количество туннелей
	tunnels := client.ListTunnels()
	if len(tunnels) != 2 {
		t.Errorf("Expected 2 tunnels, got %d", len(tunnels))
	}
	
	// Проверяем получение туннеля по ID
	tunnel, exists := client.GetTunnel("test-tunnel-1")
	if !exists {
		t.Error("Tunnel test-tunnel-1 not found")
	}
	if tunnel.ID != "test-tunnel-1" {
		t.Errorf("Expected tunnel ID 'test-tunnel-1', got '%s'", tunnel.ID)
	}
	
	// Закрываем туннель
	err := client.CloseTunnel("test-tunnel-1")
	if err != nil {
		t.Errorf("Failed to close tunnel: %v", err)
	}
	
	// Проверяем, что туннель закрыт
	tunnels = client.ListTunnels()
	if len(tunnels) != 1 {
		t.Errorf("Expected 1 tunnel after closing, got %d", len(tunnels))
	}
	
	// Проверяем, что закрытый туннель не найден
	_, exists = client.GetTunnel("test-tunnel-1")
	if exists {
		t.Error("Closed tunnel still exists")
	}
}

func TestTunnelProtocols(t *testing.T) {
	client := NewClient(false, nil)
	
	testCases := []struct {
		name     string
		protocol string
		info     map[string]interface{}
		wantErr  bool
	}{
		{
			name:     "RDP Protocol",
			protocol: "rdp",
			info: map[string]interface{}{
				"tunnel_id":    "test-rdp",
				"local_port":   3389,
				"remote_host":  "test-server",
				"remote_port":  3389,
				"protocol":     "rdp",
				"options": map[string]interface{}{
					"username": "test",
					"password": "test",
					"domain":   "test",
				},
			},
			wantErr: false,
		},
		{
			name:     "SSH Protocol",
			protocol: "ssh",
			info: map[string]interface{}{
				"tunnel_id":    "test-ssh",
				"local_port":   22,
				"remote_host":  "test-server",
				"remote_port":  22,
				"protocol":     "ssh",
				"options": map[string]interface{}{
					"username": "test",
				},
			},
			wantErr: false,
		},
		{
			name:     "Invalid Protocol",
			protocol: "invalid",
			info: map[string]interface{}{
				"tunnel_id":    "test-invalid",
				"local_port":   80,
				"remote_host":  "test-server",
				"remote_port":  80,
				"protocol":     "invalid",
				"options":      map[string]interface{}{},
			},
			wantErr: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.CreateTunnel(tc.info)
			if (err != nil) != tc.wantErr {
				t.Errorf("CreateTunnel() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestTunnelConcurrency(t *testing.T) {
	client := NewClient(false, nil)
	
	// Создаем туннели в разных горутинах
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			tunnelInfo := map[string]interface{}{
				"tunnel_id":    fmt.Sprintf("test-tunnel-%d", id),
				"local_port":   3389 + id,
				"remote_host":  "test-server",
				"remote_port":  3389,
				"protocol":     "rdp",
				"options": map[string]interface{}{
					"username": "test",
					"password": "test",
					"domain":   "test",
				},
			}
			
			_, err := client.CreateTunnel(tunnelInfo)
			if err != nil {
				t.Errorf("Failed to create tunnel %d: %v", id, err)
			}
			done <- true
		}(i)
	}
	
	// Ждем завершения всех горутин
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Проверяем количество туннелей
	tunnels := client.ListTunnels()
	if len(tunnels) != 10 {
		t.Errorf("Expected 10 tunnels, got %d", len(tunnels))
	}
} 