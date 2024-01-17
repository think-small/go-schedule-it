package db

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go-schedule-it/internal/app/features/event"
	"go-schedule-it/internal/app/logger"
	"log/slog"
	"os"
)

type EventStreamReader struct {
	db *sqlx.DB
}

func (e EventStreamReader) Read(streamId uuid.UUID) ([]event.Event, error) {
	return nil, nil
}

func NewEventStreamReader(dbProvider, dbConnStr string) *EventStreamReader {
	if db == nil {
		database, err := sqlx.Connect(dbProvider, dbConnStr)
		if err != nil {
			slog.Error("Unable to connect to database.", slog.String(logger.INNER_ERROR, err.Error()))
			os.Exit(1)
		}

		db = database
	}

	return &EventStreamReader{
		db: db,
	}
}
