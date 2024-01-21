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

type CalendarRepository interface {
	RegisterNewMeeting(*event.Event) error
	GetMeetingEvents(uuid.UUID) ([]event.Event, error)
}

type Service struct {
	calendarRepository CalendarRepository
}

func NewMeetingService(calendarRepository CalendarRepository) *Service {
	return &Service{
		calendarRepository: calendarRepository,
	}
}

func (s *Service) GetMeeting(streamId uuid.UUID) (*Meeting, error) {
	events, err := s.calendarRepository.GetMeetingEvents(streamId)
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

	err = s.calendarRepository.RegisterNewMeeting(evt)
	if err != nil {

	}

	return nil
}
