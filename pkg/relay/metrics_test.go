package relay

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMetrics(t *testing.T) {
	// Initialize metrics
	InitMetrics()

	// Record some metrics
	RecordConnection()
	RecordTunnel()
	RecordError("test_error")
	RecordHeartbeat()
	RecordHeartbeatError()
	RecordHeartbeatLatency(100 * time.Millisecond)

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(MetricsHandler))
	defer server.Close()

	// Test metrics endpoint
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to get metrics: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
}

func TestHealthCheck(t *testing.T) {
	// Initialize metrics
	InitMetrics()

	// Record some metrics
	RecordConnection()
	RecordTunnel()
	RecordError("test_error")

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(HealthCheckHandler))
	defer server.Close()

	// Test health check endpoint
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to get health check: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}

	// Verify health status
	status := GetHealthStatus()
	if status.Status != "ok" {
		t.Errorf("Expected status 'ok', got %v", status.Status)
	}
}

func TestMetricsConcurrency(t *testing.T) {
	// Initialize metrics
	InitMetrics()

	// Test concurrent access
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			RecordConnection()
			RecordTunnel()
			RecordError("test_error")
			RecordHeartbeat()
			RecordHeartbeatError()
			RecordHeartbeatLatency(100 * time.Millisecond)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestMetricsPersistence(t *testing.T) {
	// Initialize metrics
	InitMetrics()

	// Record initial metrics
	RecordConnection()
	RecordTunnel()
	RecordError("test_error")

	// Get initial health status
	initialStatus := GetHealthStatus()

	// Record more metrics
	RecordConnection()
	RecordTunnel()
	RecordError("test_error")

	// Get updated health status
	updatedStatus := GetHealthStatus()

	// Verify metrics are persisted
	if initialStatus.Status != updatedStatus.Status {
		t.Errorf("Expected status to be consistent, got %v and %v", initialStatus.Status, updatedStatus.Status)
	}
} 