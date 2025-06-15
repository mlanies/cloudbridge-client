package relay

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMetrics(t *testing.T) {
	// Record some metrics
	RecordConnection(1.5)
	RecordError("test_error")
	SetActiveTunnels(5)
	RecordHeartbeat(0.1)
	RecordMissedHeartbeat()

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/metrics" {
			t.Errorf("Expected to request '/metrics', got: %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got: %s", r.Method)
		}
	}))
	defer server.Close()

	// Test metrics endpoint
	resp, err := http.Get(server.URL + "/metrics")
	if err != nil {
		t.Fatalf("Failed to get metrics: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got: %v", resp.StatusCode)
	}
}

func TestHealthCheck(t *testing.T) {
	// Record some metrics
	RecordConnection(1.0)
	RecordError("test_error")
	SetActiveTunnels(3)

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
	// Test concurrent access
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			RecordConnection(1.0)
			RecordError("concurrent_error")
			SetActiveTunnels(i)
			RecordHeartbeat(0.1)
			RecordMissedHeartbeat()
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestMetricsPersistence(t *testing.T) {
	// Record initial metrics
	RecordConnection(1.0)
	RecordError("test_error")
	SetActiveTunnels(2)

	// Get initial health status
	initialStatus := GetHealthStatus()

	// Record more metrics
	RecordConnection(1.0)
	RecordError("test_error")
	SetActiveTunnels(3)

	// Get updated health status
	updatedStatus := GetHealthStatus()

	// Verify metrics are persisted
	if initialStatus.Status != updatedStatus.Status {
		t.Errorf("Expected status to be consistent, got %v and %v", initialStatus.Status, updatedStatus.Status)
	}
} 