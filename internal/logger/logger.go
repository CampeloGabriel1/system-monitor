package logger

import (
	"log/slog"
	"os"
)

func SetupLogger() *slog.Logger {
	// Handler para o Console (JSON)
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)
	return logger
}