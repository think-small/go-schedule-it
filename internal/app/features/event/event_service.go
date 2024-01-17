package event

import (
	"github.com/google/uuid"
)

type StreamReader interface {
	Read(uuid uuid.UUID) ([]Event, error)
}

type StreamWriter interface {
	Write(Event) error
}

type Service struct {
	writer *StreamWriter
	reader *StreamReader
}

func (s Service) createCalendarEvent(e Event) error {
	return nil
}

func (s Service) cancelCalendarEvent(e Event) error {
	return nil
}

func (s Service) startCalendarEvent(e Event) error {
	return nil
}

func (s Service) endCalendarEvent(e Event) error {
	return nil
}

func NewEventService(writer StreamWriter, reader StreamReader) *Service {
	return &Service{
		writer: &writer,
		reader: &reader,
	}
}
