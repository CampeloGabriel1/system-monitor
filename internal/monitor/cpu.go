package monitor

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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
	found := false

	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") {
			f := strings.Fields(line)
			if len(f) < 5 {
				return CPUStats{}, fmt.Errorf("formato inválido de /proc/stat: esperado 5 campos")
			}

			var errParse error
			if c.User, errParse = strconv.ParseUint(f[1], 10, 64); errParse != nil {
				return CPUStats{}, fmt.Errorf("erro ao fazer parse de User: %v", errParse)
			}
			if c.Nice, errParse = strconv.ParseUint(f[2], 10, 64); errParse != nil {
				return CPUStats{}, fmt.Errorf("erro ao fazer parse de Nice: %v", errParse)
			}
			if c.System, errParse = strconv.ParseUint(f[3], 10, 64); errParse != nil {
				return CPUStats{}, fmt.Errorf("erro ao fazer parse de System: %v", errParse)
			}
			if c.Idle, errParse = strconv.ParseUint(f[4], 10, 64); errParse != nil {
				return CPUStats{}, fmt.Errorf("erro ao fazer parse de Idle: %v", errParse)
			}

			c.Total = c.User + c.Nice + c.System + c.Idle
			found = true
			break
		}
	}

	if !found {
		return CPUStats{}, fmt.Errorf("linha 'cpu' não encontrada em /proc/stat")
	}

	return c, nil
}