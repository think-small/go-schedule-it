package main

import (
	"github.com/subosito/gotenv"
	"go-schedule-it/internal/app/logger"
	"go-schedule-it/internal/app/server"
	"log/slog"
	"os"
)

func main() {
	err := gotenv.Load()
	if err != nil {
		slog.Error("Unable to load config.", slog.String(logger.INNER_ERROR, err.Error()))
		os.Exit(1)
	}

	server.Run()
}
