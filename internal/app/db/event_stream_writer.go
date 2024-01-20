package db

import (
	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go-schedule-it/internal/app/features/event"
	"go-schedule-it/internal/app/logger"
	"log/slog"
	"os"
)

type EventStreamWriter struct {
	db *sqlx.DB
}

func (e EventStreamWriter) Write(evt event.Event) error {
	sql := `
		INSERT INTO calendar_events (streamid, eventtype, payload)
		VALUES ( $1, $2, $3)`

	_, err := db.Exec(sql, evt.StreamId, evt.EventType, evt.Payload)
	if err != nil {
		slog.Error("Unable to persist event to database.", slog.String(logger.INNER_ERROR, err.Error()))
		return err
	}
	return nil
}

func NewEventStreamWriter(dbProvider, dbConnStr string) *EventStreamWriter {
	if db == nil {
		database, err := sqlx.Connect(dbProvider, dbConnStr)
		if err != nil {
			slog.Error("Unable to connect to database.", slog.String(logger.INNER_ERROR, err.Error()))
			os.Exit(1)
		}

		db = database
	}

	return &EventStreamWriter{
		db: db,
	}
}
