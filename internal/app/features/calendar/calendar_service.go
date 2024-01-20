package calendar

import (
	"github.com/google/uuid"
	"go-schedule-it/internal/app/features/event"
)

type Service struct {
	reader event.StreamReader
}

func NewCalendarService(sr event.StreamReader) *Service {
	return &Service{
		reader: sr,
	}
}

func (s *Service) GetCalendar(streamId uuid.UUID) (*Calendar, error) {
	events, err := s.reader.Read(streamId)
	if err != nil {
		return nil, err
	}

	calendar := &Calendar{}
	err = calendar.Apply(events)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
