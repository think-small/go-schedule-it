package db

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go-schedule-it/internal/app/logger"
	"log/slog"
)

type RelReader struct {
	db *sqlx.DB
}

func (r RelReader) Read(calendarId uuid.UUID) ([]uuid.UUID, error) {
	var meetings []uuid.UUID

	sql := `
		SELECT
			calendarId,
			meetingId
		FROM calendar_meetings
		WHERE calendarId = $1
	`
	err := r.db.Select(meetings, sql, calendarId)
	if err != nil {
		slog.Error("Unable to get calendar meeting joins.", slog.String(logger.INNER_ERROR, err.Error()))
		return nil, err
	}

	return meetings, nil
}
