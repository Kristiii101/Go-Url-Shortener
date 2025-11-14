package observability

import (
	"log"
	"os"
)

func NewLogger() *log.Logger {
	// Simple stdout logger; you can swap for slog/zap later.
	return log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
}
