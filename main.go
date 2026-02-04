package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
)

	
func main () {
	
	content, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler o arquivo: %v", err)
		os.Exit(1)
	}
	texto := string(content)
	linhas := strings.Split(texto, "\n")

	var memTotal, memAvailable int
	for _, linha := range linhas {
		if strings.HasPrefix(linha, "MemTotal:") {

			fields := strings.Fields(linha)
			if len(fields) > 1 {
				if v, err := strconv.Atoi(fields[1]); err == nil {
					memTotal = v
				} else {
					fmt.Fprintf(os.Stderr, "Erro ao converter valor de MemTotal: %v\n", err)
				}
			}
		} else if strings.HasPrefix(linha, "MemAvailable:") {
			fields := strings.Fields(linha)
			if len(fields) > 1 {
				if v, err := strconv.Atoi(fields[1]); err == nil {
					memAvailable = v
				} else {
					fmt.Fprintf(os.Stderr, "Erro ao converter valor de MemTotal: %v\n", err)
				}
			}
		}
	}
	fmt.Printf("MemTotal (kB): %d\nMemAvailable (kB): %d\n", memTotal, memAvailable)

}