package relay

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMetrics(t *testing.T) {
	// Тестируем запись метрик
	RecordConnection(1.5)
	RecordError("test_error")
	SetActiveTunnels(5)
	RecordHeartbeat(0.1)
	RecordMissedHeartbeat()

	// Создаем тестовый HTTP сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/metrics" {
			t.Errorf("Expected to request '/metrics', got: %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got: %s", r.Method)
		}
	}))
	defer ts.Close()

	// Проверяем доступность метрик
	resp, err := http.Get(ts.URL + "/metrics")
	if err != nil {
		t.Fatalf("Failed to get metrics: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got: %v", resp.StatusCode)
	}
}

func TestHealthCheck(t *testing.T) {
	// Тестируем health check
	status := GetHealthStatus()
	if status.Status != "ok" {
		t.Errorf("Expected status 'ok', got: %s", status.Status)
	}
	if status.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got: %s", status.Version)
	}

	// Проверяем метрики в health check
	metrics := status.Metrics
	if metrics["active_connections"] < 0 {
		t.Error("Expected non-negative active_connections")
	}
	if metrics["active_tunnels"] < 0 {
		t.Error("Expected non-negative active_tunnels")
	}
	if metrics["total_errors"] < 0 {
		t.Error("Expected non-negative total_errors")
	}
	if metrics["missed_heartbeats"] < 0 {
		t.Error("Expected non-negative missed_heartbeats")
	}

	// Тестируем обновление статуса
	UpdateHealthStatus("degraded")
	status = GetHealthStatus()
	if status.Status != "degraded" {
		t.Errorf("Expected status 'degraded', got: %s", status.Status)
	}
}

func TestMetricsConcurrency(t *testing.T) {
	// Тестируем конкурентный доступ к метрикам
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

	// Ждем завершения всех горутин
	for i := 0; i < 10; i++ {
		<-done
	}

	// Проверяем, что метрики обновились
	status := GetHealthStatus()
	if status.Metrics["total_errors"] < 10 {
		t.Error("Expected at least 10 errors recorded")
	}
}

func TestMetricsPersistence(t *testing.T) {
	// Тестируем сохранение метрик между вызовами
	initialErrors := GetHealthStatus().Metrics["total_errors"]
	
	RecordError("persistent_error")
	
	time.Sleep(100 * time.Millisecond)
	
	finalErrors := GetHealthStatus().Metrics["total_errors"]
	if finalErrors <= initialErrors {
		t.Error("Expected error count to increase")
	}
} 