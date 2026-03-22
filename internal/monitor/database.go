package monitor

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MetricLog struct {
	ID              uint      `gorm:"primaryKey"`
	Timestamp       time.Time `gorm:"index"`
	CpuUsage        float64
	MemoryUsage     float64
	TotalMemoryKb   int
}

var DB *gorm.DB

func InitDB() {
	// Pega a URL do banco que definimos no docker-compose.yml
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		// Fallback para caso rode fora do docker (localhost)
		dsn = "host=localhost user=user password=password dbname=monitor_db port=5432 sslmode=disable"
	}

	var err error
    // Tenta conectar 5 vezes com pausa de 2 segundos entre elas
    for i := 0; i < 5; i++ {
        DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {
            break
        }
        log.Printf("Aguardando banco de dados... (tentativa %d/5)", i+1)
        time.Sleep(2 * time.Second)
    }

    if err != nil {
        log.Fatalf("Falha final ao conectar no banco: %v", err)
    }

    DB.AutoMigrate(&MetricLog{})
    log.Println("✅ Banco de dados conectado!")
}