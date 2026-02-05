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
			corMem := getCor(usoMem)
			fmt.Printf("RAM: %s%.2f%%%s de %d kB usados\n", corMem, usoMem, Reset, mStats.Total)
		}

		// 2. CPU
		statsA, err := GetCPUStats()
		if err != nil {
			fmt.Printf("Erro CPU (1ª leitura): %v\n", err)
			continue
		}

		time.Sleep(1 * time.Second)

		statsB, err := GetCPUStats()
		if err != nil {
			fmt.Printf("Erro CPU (2ª leitura): %v\n", err)
			continue
		}

		idleDelta := float64(statsB.Idle - statsA.Idle)
		totalDelta := float64(statsB.Total - statsA.Total)

		if totalDelta == 0 {
			fmt.Println("Erro: Delta total de CPU é zero")
			continue
		}

		cpuUsage := (1.0 - idleDelta/totalDelta) * 100

		corCPU := getCor(cpuUsage)
		fmt.Printf("CPU: %s%.2f%%%s\n", corCPU, cpuUsage, Reset)
	}
}