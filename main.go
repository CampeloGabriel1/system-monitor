package main

import (
	"fmt"
	"time"
)

func main() {
	// 1. Memória
	mStats, err := GetMemoryStats()
	if err != nil {
		fmt.Printf("Erro memória: %v\n", err)
	} else {
		fmt.Printf("RAM: %.2f%% de %d kB usados\n", mStats.UsagePercentage(), mStats.Total)
	}

	// 2. CPU (Cálculo de Delta)
	statsA, _ := GetCPUStats()
	time.Sleep(1 * time.Second)
	statsB, _ := GetCPUStats()

	idleDelta := float64(statsB.Idle - statsA.Idle)
	totalDelta := float64(statsB.Total - statsA.Total)
	cpuUsage := (1.0 - idleDelta/totalDelta) * 100

	fmt.Printf("CPU: %.2f%%\n", cpuUsage)
}