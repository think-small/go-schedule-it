package logger

import (
	"log/slog"
	"os"
)

const (
	INNER_ERROR = "Inner Error"
	COUNT       = "Count"
)

func init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
