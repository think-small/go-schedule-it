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

func (s Service) createMeeting(e Event) error {
	return nil
}

func (s Service) cancelMeeting(e Event) error {
	return nil
}

func (s Service) startMeeting(e Event) error {
	return nil
}

func (s Service) endMeeting(e Event) error {
	return nil
}

func NewEventService(writer StreamWriter, reader StreamReader) *Service {
	return &Service{
		writer: &writer,
		reader: &reader,
	}
}
