package db

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go-schedule-it/internal/app/logger"
	"log/slog"
	"os"
)

type RelWriter struct {
	db *sqlx.DB
}

func (r RelWriter) Write(calendarId uuid.UUID, meetingId uuid.UUID) error {
	sql := `
		INSERT INTO calendar_meetings (calendarId, meetingId)
		VALUES($1, $2)
	`

	_, err := r.db.Exec(sql, calendarId, meetingId)
	if err != nil {
		slog.Error("Unable to add meeting to calendar.", slog.String(logger.INNER_ERROR, err.Error()))
		return err
	}

	return nil
}

func NewRelWriter(dbProvider, dbConnStr string) *RelWriter {
	if db == nil {
		database, err := sqlx.Connect(dbProvider, dbConnStr)
		if err != nil {
			slog.Error("Unable to connect to database.", slog.String(logger.INNER_ERROR, err.Error()))
			os.Exit(1)
		}

		db = database
	}

	return &RelWriter{
		db: db,
	}
}
