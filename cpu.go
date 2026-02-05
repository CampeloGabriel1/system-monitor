package main

import (
	"os"
	"strings"
	"strconv"
)

type CPUStats struct {
	User, Nice, System, Idle, Total uint64
}

func GetCPUStats() (CPUStats, error) {
	content, err := os.ReadFile("/proc/stat")
	if err != nil {
		return CPUStats{}, err
	}

	lines := strings.Split(string(content), "\n")
	var c CPUStats

	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") {
			f := strings.Fields(line)
			if len(f) > 4 {
				c.User, _ = strconv.ParseUint(f[1], 10, 64)
				c.Nice, _ = strconv.ParseUint(f[2], 10, 64)
				c.System, _ = strconv.ParseUint(f[3], 10, 64)
				c.Idle, _ = strconv.ParseUint(f[4], 10, 64)
				c.Total = c.User + c.Nice + c.System + c.Idle
			}
			break
		}
	}
	return c, nil
}