package main

import (
	"fmt"
	"os"
)

	
func main () {
	content, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler o arquivo: %v", err)
		os.Exit(1)
	}
	fmt.Println(string(content))
}