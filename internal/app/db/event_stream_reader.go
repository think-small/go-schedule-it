package db

import (
	"github.com/google/uuid"
	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/v5/stdlib"
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
	var events []event.Event
	sql := `
		SELECT
			id,
			streamId,
			eventType,
			version,
			timestamp,
			payload
		FROM events.calendar_events
		WHERE streamId = $1
		ORDER BY version
	`

	err := db.Select(events, sql)
	if err != nil {
		slog.Error("Unable to retrieve events from the database.", slog.String(logger.INNER_ERROR, err.Error()))
		return nil, err
	}

	return events, nil
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
