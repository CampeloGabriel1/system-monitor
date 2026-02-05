package main

import (
	"os"
	"strings"
	"strconv"
)

type MemoryStats struct {
	Total     int
	Available int
	Used      int
}

func (m MemoryStats) UsagePercentage() float64 {
	if m.Total == 0 { return 0 }
	return float64(m.Used) / float64(m.Total) * 100
}

func GetMemoryStats() (MemoryStats, error) {
	content, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return MemoryStats{}, err
	}

	lines := strings.Split(string(content), "\n")
	var m MemoryStats

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 { continue }

		val, _ := strconv.Atoi(fields[1])
		if strings.HasPrefix(line, "MemTotal:") {
			m.Total = val
		} else if strings.HasPrefix(line, "MemAvailable:") {
			m.Available = val
		}
	}
	m.Used = m.Total - m.Available
	return m, nil
}