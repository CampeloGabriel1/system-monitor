package main

// teste
import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"system-monitor/internal/api"
	"system-monitor/internal/logger"
	"system-monitor/internal/monitor"
)

func main() {
	
	monitor.InitDB()
	http.HandleFunc("/", api.HomeHandler)
	http.HandleFunc("/health", api.HealthHandler)
	http.HandleFunc("/history", api.HistoryHandler)
	http.HandleFunc("/stats", api.StatsHandler)
	http.HandleFunc("/memory", api.MemoryHandler)
	http.HandleFunc("/cpu", api.CPUHandler)

	// Porta: usa variável de ambiente PORT se existir, senão 8080 (ex.: PORT=3000 ./monitor).
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if port[0] != ':' {
		port = ":" + port
	}
	fmt.Printf("Servidor de Monitor iniciado em http://localhost%s\n", port)
	fmt.Println("Endpoints disponíveis:")
	fmt.Println("  GET /health  - Verifica saúde do servidor")
	fmt.Println("  GET /history - Histórico de métricas salvas")
	fmt.Println("  GET /stats   - Stats de CPU e Memória")
	fmt.Println("  GET /memory  - Stats de Memória")
	fmt.Println("  GET /cpu     - Stats de CPU")

	logg := logger.SetupLogger()
	slog.SetDefault(logg)

	slog.Info("Iniciando serviço de monitoramento de alta escala")

	go func() {
		for {
			stats := monitor.GetStats()

			if stats.Status == "WARNING" {
				slog.Warn("Uso de memória acima do limite!", 
					"mem_kb", stats.MemoryAlloc, 
					"status", stats.Status)
			} else {
				slog.Info("Sistema Saudável", 
					"mem_kb", stats.MemoryAlloc)
			}

			time.Sleep(10 * time.Second)
		}
	}()

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}

}