package event

import (
	"encoding/json"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type EventVM struct {
	StreamId  uuid.UUID `json:"streamId"`
	EventType EventType `json:"eventType"`
	Payload   json.RawMessage
}

func (e *EventVM) IsValid() bool {
	if e.StreamId == uuid.Nil {
		slog.Info("No streamId provided.")
		return false
	}
	if e.EventType.IsValid() == false {
		slog.Info("Invalid event type provided.")
		return false
	}
	if e.Payload == nil {
		slog.Info("No event payload provided.")
		return false
	}

	return true
}

type Event struct {
	Id        uuid.UUID       `json:"id"`
	StreamId  uuid.UUID       `json:"streamId"`
	EventType EventType       `json:"eventType"`
	Version   int8            `json:"version"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

type EventType int

const (
	Unsupported EventType = iota
	Created
	Canceled
	Started
	Ended
	AttendeeRegistered
	AttendeeUnregistered
)

func (et EventType) IsValid() bool {
	return et > 0 && et < Ended
}

type CalendarEventCreated struct {
	ScheduledStart time.Time `json:"scheduledStart"`
	ScheduledEnd   time.Time `json:"scheduledEnd"`
	HostId         uuid.UUID `json:"hostId"`
}

func (c CalendarEventCreated) IsValid() bool {
	if c.HostId == uuid.Nil {
		slog.Info("No hostId provided.")
		return false
	}
	if c.ScheduledStart.IsZero() {
		slog.Info("No scheduledStart provided.")
		return false
	}
	if c.ScheduledEnd.IsZero() {
		slog.Info("No scheduledEnd provided.")
		return false
	}
	if c.ScheduledStart.Before(time.Now().UTC()) {
		slog.Info("Provided scheduledStart is in the past.")
		return false
	}

	return true
}

type CalendarEventCanceled struct {
	EventId    uuid.UUID `json:"evnetId"`
	CanceledAt time.Time `json:"canceledAt"`
}

type CalendarEventStarted struct {
	EventId     uuid.UUID `json:"eventId"`
	ActualStart time.Time `json:"actualStart"`
}

type CalendarEventEnded struct {
	EventId   uuid.UUID `json:"eventId"`
	ActualEnd time.Time `json:"actualEnd"`
}

type CalendarEventAttendeeRegistered struct {
	EventId    uuid.UUID `json:"eventId"`
	AttendeeId uuid.UUID `json:"attendeeId"`
}

type CalendarEventAttendeeUnregistered struct {
	EventId    uuid.UUID `json:"eventId"`
	AttendeeId uuid.UUID `json:"attendeeId"`
}
