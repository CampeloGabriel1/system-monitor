package main

import (
	"fmt"
	"time"
)

const (
    Reset  = "\033[0m"
    Red    = "\033[31m"
    Yellow = "\033[33m"
    Green  = "\033[32m"
    Cyan   = "\033[36m"
	Clear  = "\033[H\033[2J"
)

func getCor(uso float64) string {
	if uso > 80 {
		return Red
	}
	if uso > 50 {
		return Yellow
	}
	return Green
}

func main() {
	for {
	fmt.Print(Clear)

		fmt.Println(Cyan + "=== GO SYSTEM MONITOR (Pressione Ctrl+C para sair) ===" + Reset)
		fmt.Println(time.Now().Format("15:04:05"))
		fmt.Println("--------------------------------------------------")

	// 1. Memória
	mStats, err := GetMemoryStats()
	if err != nil {
		fmt.Printf("Erro memória: %v\n", err)
	} else {
		usoMem := mStats.UsagePercentage()
        corMem := getCor(usoMem) // Aplicando a cor na RAM também!
        fmt.Printf("RAM: %s%.2f%%%s de %d kB usados\n", corMem, usoMem, Reset, mStats.Total)
	}

	// 2. CPU 
	statsA, _ := GetCPUStats()
	time.Sleep(1 * time.Second)
	statsB, _ := GetCPUStats()

	idleDelta := float64(statsB.Idle - statsA.Idle)
	totalDelta := float64(statsB.Total - statsA.Total)
	cpuUsage := (1.0 - idleDelta/totalDelta) * 100

	corCPU := getCor(cpuUsage)
	fmt.Printf("CPU: %s%.2f%%%s\n", corCPU, cpuUsage, Reset)
	}
}