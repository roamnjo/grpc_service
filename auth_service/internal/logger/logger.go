package logger

import (
	"log/slog"
	"os"
)

func New(level slog.Leveler) *slog.Logger {
	opts := &slog.HandlerOptions{Level: level}
	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}
