package main

import (
	"fmt"
	"os"
	"strings"
)

	
func main () {
	
	content, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler o arquivo: %v", err)
		os.Exit(1)
	}
	texto := string(content)
	linhas := strings.Split(texto, "\n")
	var resultado []string
	for _, linha := range linhas {
		if strings.HasPrefix(linha, "MemTotal:") {
			resultado = append(resultado, linha)
		} else if strings.HasPrefix(linha, "MemAvailable:") {
			resultado = append(resultado, linha)
		}
	}
	for _, linha := range resultado {
		fmt.Println(linha)
	}

}