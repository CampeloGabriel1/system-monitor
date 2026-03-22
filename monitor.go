package main

import (
	"runtime"
	"time"
)

type SystemInfo struct {
	Timestamp   string  `json:"timestamp"`
	MemoryAlloc uint64  `json:"memory_alloc_kb"`
	Status      string  `json:"status"`
}

func GetStats() SystemInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	status := "OK"
	if m.Alloc/1024 > 81920 { // 80MB Threshold
		status = "WARNING"
	}

	return SystemInfo{
		Timestamp:   time.Now().Format(time.RFC3339),
		MemoryAlloc: m.Alloc / 1024,
		Status:      status,
	}
}