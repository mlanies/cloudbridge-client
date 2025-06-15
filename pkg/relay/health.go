package relay

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthStatus представляет текущее состояние сервера
type HealthStatus struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Uptime    string           `json:"uptime"`
	StartTime time.Time        `json:"start_time"`
	Metrics   map[string]int64 `json:"metrics"`
}

var (
	startTime = time.Now()
	healthStatus = &HealthStatus{
		Status:    "ok",
		Version:   "1.0.0",
		StartTime: startTime,
		Metrics:   make(map[string]int64),
	}
)

// UpdateHealthStatus обновляет статус здоровья
func UpdateHealthStatus(status string) {
	healthStatus.Status = status
}

// GetHealthStatus возвращает текущий статус здоровья
func GetHealthStatus() *HealthStatus {
	healthStatus.Uptime = time.Since(startTime).String()
	
	// Обновляем метрики
	healthStatus.Metrics = map[string]int64{
		"active_connections": int64(connectionsTotal),
		"active_tunnels":     int64(activeTunnels),
		"total_errors":       int64(errorsTotal),
		"missed_heartbeats":  int64(missedHeartbeats),
	}
	
	return healthStatus
}

// HealthCheckHandler обрабатывает запросы к /health
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := GetHealthStatus()
	
	w.Header().Set("Content-Type", "application/json")
	if status.Status != "ok" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	
	json.NewEncoder(w).Encode(status)
} 