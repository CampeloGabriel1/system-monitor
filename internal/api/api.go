package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"system-monitor/internal/monitor"
)

type StatsResponse struct {
	Timestamp string            `json:"timestamp"`
	Memory    MemoryStatsJSON   `json:"memory"`
	CPU       CPUStatsJSON      `json:"cpu"`
}

type MemoryStatsJSON struct {
	Total            int     `json:"total_kb"`
	Available        int     `json:"available_kb"`
	Used             int     `json:"used_kb"`
	UsagePercentage  float64 `json:"usage_percentage"`
}

type CPUStatsJSON struct {
	User            uint64  `json:"user"`
	Nice            uint64  `json:"nice"`
	System          uint64  `json:"system"`
	Idle            uint64  `json:"idle"`
	Total           uint64  `json:"total"`
	UsagePercentage float64 `json:"usage_percentage"`
}

// resultMemory e resultCPU carregam o retorno das goroutines (dados + erro).
// Usamos structs para enviar os dois valores pelo canal de uma vez.
type resultMemory struct {
	stats monitor.MemoryStats
	err   error
}
type resultCPU struct {
	stats monitor.CPUStats
	err   error
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"project":     "Sentinel Core - System Monitor",
		"status":      "running",
		"description": "API para monitoramento de recursos Linux em tempo real e histórico.",
		"endpoints": []string{
			"GET /health   - Verifica saúde do sistema",
			"GET /stats    - Métricas atuais de CPU e Memória (salva no banco)",
			"GET /memory   - Somente métricas de memória",
			"GET /cpu      - Somente métricas de CPU",
			"GET /history  - Histórico das últimas 20 métricas salvas",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("home: erro ao escrever resposta JSON: %v", err)
	}
}

func HistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var metrics []monitor.MetricLog

	result := monitor.DB.Order("timestamp desc").Limit(20).Find(&metrics)
	if result.Error != nil {
		log.Printf("Erro ao buscar histórico de métricas: %v", result.Error)
		http.Error(w, "Erro ao processar o histórico de métricas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		log.Printf("history: erro ao escrever resposta JSON: %v", err)
	}
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Canais onde as goroutines vão enviar o resultado (buffer 1 = não bloqueia ao enviar).
	memCh := make(chan resultMemory, 1)
	cpuCh := make(chan resultCPU, 1)

	// Goroutine 1: busca memória em paralelo.
	go func() {
		stats, err := monitor.GetMemoryStats()
		memCh <- resultMemory{stats, err}
	}()

	// Goroutine 2: busca CPU em paralelo.
	go func() {
		stats, err := monitor.GetCPUStats()
		cpuCh <- resultCPU{stats, err}
	}()

	// Esperar os dois resultados (a ordem não importa; quem chegar primeiro espera o outro).
	memRes := <-memCh
	cpuRes := <-cpuCh

	if memRes.err != nil {
		http.Error(w, "Error getting memory stats", http.StatusInternalServerError)
		return
	}
	if cpuRes.err != nil {
		http.Error(w, "Error getting CPU stats", http.StatusInternalServerError)
		return
	}

	memstats := memRes.stats
	cStats := cpuRes.stats
	cpuUsage := calculateCPUUsage(cStats)

	response := StatsResponse{
		Timestamp: time.Now().Format(time.RFC3339),
		Memory: MemoryStatsJSON{
			Total:           memstats.Total,
			Available:       memstats.Available,
			Used:            memstats.Used,
			UsagePercentage: memstats.UsagePercentage(),
		},
		CPU: CPUStatsJSON{
			User:            cStats.User,
			Nice:            cStats.Nice,
			System:          cStats.System,
			Idle:            cStats.Idle,
			Total:           cStats.Total,
			UsagePercentage: cpuUsage,
		},
	}

	// Salvar métrica no banco de dados
	err := SaveMetric(cpuUsage, memstats.UsagePercentage(), memstats.Total)
	if err != nil {
		log.Printf("Erro ao salvar métrica: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("stats: erro ao escrever resposta JSON: %v", err)
	}
}

func MemoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	memstats, err := monitor.GetMemoryStats()
	if err != nil {
		http.Error(w, "Error getting memory stats", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"timestamp":       time.Now().Format(time.RFC3339),
		"total_kb":        memstats.Total,
		"available_kb":    memstats.Available,
		"used_kb":         memstats.Used,
		"usage_percentage": memstats.UsagePercentage(),
	}

	// Salvar métrica no banco de dados
	err = SaveMetric(0, memstats.UsagePercentage(), memstats.Total)
	if err != nil {
		log.Printf("Erro ao salvar métrica: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("memory: erro ao escrever resposta JSON: %v", err)
	}
}

func CPUHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cStats, err := monitor.GetCPUStats()
	if err != nil {
		http.Error(w, "Error getting CPU stats", http.StatusInternalServerError)
		return
	}

	cpuUsage := calculateCPUUsage(cStats)

	response := map[string]interface{}{
		"timestamp":        time.Now().Format(time.RFC3339),
		"user":             cStats.User,
		"nice":             cStats.Nice,
		"system":           cStats.System,
		"idle":             cStats.Idle,
		"total":            cStats.Total,
		"usage_percentage": cpuUsage,
	}

	// Salvar métrica no banco de dados
	err = SaveMetric(cpuUsage, 0, 0)
	if err != nil {
		log.Printf("Erro ao salvar métrica: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("cpu: erro ao escrever resposta JSON: %v", err)
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("health: erro ao escrever resposta JSON: %v", err)
	}
}

// SaveMetric salva uma métrica no banco de dados
func SaveMetric(cpuUsage float64, memoryUsage float64, totalMemoryKb int) error {
	metric := monitor.MetricLog{
		Timestamp:     time.Now(),
		CpuUsage:      cpuUsage,
		MemoryUsage:   memoryUsage,
		TotalMemoryKb: totalMemoryKb,
	}
	return monitor.DB.Create(&metric).Error
}

func calculateCPUUsage(stats monitor.CPUStats) float64 {
	if stats.Total == 0 {
		return 0
	}
	used := stats.User + stats.Nice + stats.System
	return float64(used) / float64(stats.Total) * 100
}