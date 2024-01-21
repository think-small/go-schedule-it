package db

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go-schedule-it/internal/app/features/event"
	"go-schedule-it/internal/app/features/meeting"
	"go-schedule-it/internal/app/logger"
	"log/slog"
	"os"
)

type CalendarRepository struct {
	db        *sqlx.DB
	relWriter *RelWriter
	relReader *RelReader
	evtWriter *EventStreamWriter
	evtReader *EventStreamReader
}

func NewCalendarRepository(dbProvider, dbConnStr string, relWriter *RelWriter, relReader *RelReader, evtWriter *EventStreamWriter, evtReader *EventStreamReader) *CalendarRepository {
	if db == nil {
		database, err := sqlx.Connect(dbProvider, dbConnStr)
		if err != nil {
			slog.Error("Unable to connect to database.", slog.String(logger.INNER_ERROR, err.Error()))
			os.Exit(1)
		}

		db = database
	}

	return &CalendarRepository{
		db:        db,
		relWriter: relWriter,
		relReader: relReader,
		evtWriter: evtWriter,
		evtReader: evtReader,
	}
}

func (c CalendarRepository) GetMeetingEvents(streamId uuid.UUID) ([]event.Event, error) {
	events, err := c.evtReader.Read(streamId)
	if err != nil {
		slog.Error("Unable to retrieve events for meeting.", slog.String(logger.INNER_ERROR, err.Error()))
		return nil, err
	}

	return events, nil
}

func (c CalendarRepository) RegisterNewMeeting(evt *event.Event) error {
	var err error
	var eventPayload = &event.MeetingCreated{}

	if evt.EventType != event.Created {
		err = &meeting.InvalidEventTypeError{}
		slog.Info(err.Error())
		return err
	}

	err = json.Unmarshal(evt.Payload, eventPayload)
	if err != nil {
		slog.Error("Unable to unpack event payload. Unable to register new meeting.")
		return err
	}

	if eventPayload.IsValid() == false {
		slog.Error("Invalid event payload provided. Unable to register new meeting.")
		return errors.New("invalid event payload provided. Unable to register new meeting")
	}

	tx, err := db.Begin()
	if err != nil {
		slog.Error("Unable to establish transaction.", slog.String(logger.INNER_ERROR, err.Error()))
		err = tx.Rollback()
		if err != nil {
			slog.Error("Unable to rollback transaction when registering meeting.", logger.INNER_ERROR, err.Error())
			return err
		}
		return err
	}

	err = c.relWriter.Write(eventPayload.CalendarId, evt.StreamId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = c.evtWriter.Write(*evt)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {

	}

	return nil
}
