package meeting

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"go-schedule-it/internal/app/features/event"
	"log/slog"
)

type RelationshipReader interface {
	Read(calendarId uuid.UUID) ([]uuid.UUID, error)
}

type RelationshipWriter interface {
	Write(calendarId uuid.UUID, meetingId uuid.UUID) error
}

type Service struct {
	streamReader event.StreamReader
	relWriter    RelationshipWriter
}

func NewMeetingService(sr event.StreamReader, rw RelationshipWriter) *Service {
	return &Service{
		streamReader: sr,
		relWriter:    rw,
	}
}

func (s *Service) GetMeeting(streamId uuid.UUID) (*Meeting, error) {
	events, err := s.streamReader.Read(streamId)
	if err != nil {
		return nil, err
	}

	meeting := &Meeting{}
	err = meeting.Apply(events)
	if err != nil {
		return nil, err
	}

	return meeting, nil
}

func (s *Service) RegisterNewMeeting(evt *event.Event) error {
	var eventPayload = &event.MeetingCreated{}
	var err error

	if evt.EventType != event.Created {
		err = &InvalidEventTypeError{}
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

	// TOOD - BS - need to perform insertion into calendar_meetings and meeting_events in a single transaction.
	//			   need to compose a db repo to combine relWriter and eventStreamWriter
	err = s.relWriter.Write(eventPayload.CalendarId, evt.StreamId)
	if err != nil {

	}

	return nil
}
