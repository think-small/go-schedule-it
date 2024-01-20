package event

import (
	"encoding/json"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type ViewModel struct {
	StreamId  uuid.UUID `json:"streamId"`
	EventType Type      `json:"eventType"`
	Payload   json.RawMessage
}

func (e *ViewModel) IsValid() bool {
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
	EventType Type            `json:"eventType"`
	Version   int8            `json:"version"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

type Type int

const (
	Unsupported Type = iota
	Created
	Canceled
	Started
	Ended
	AttendeeRegistered
	AttendeeUnregistered
)

func (et Type) IsValid() bool {
	return et > 0 && et < Ended
}

type MeetingCreated struct {
	ScheduledStart time.Time `json:"scheduledStart"`
	ScheduledEnd   time.Time `json:"scheduledEnd"`
	HostId         uuid.UUID `json:"hostId"`
}

func (c MeetingCreated) IsValid() bool {
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

type MeetingCanceled struct {
	EventId    uuid.UUID `json:"evnetId"`
	CanceledAt time.Time `json:"canceledAt"`
}

type MeetingStarted struct {
	EventId     uuid.UUID `json:"eventId"`
	ActualStart time.Time `json:"actualStart"`
}

type MeetingEnded struct {
	EventId   uuid.UUID `json:"eventId"`
	ActualEnd time.Time `json:"actualEnd"`
}

type MeetingAttendeeRegistered struct {
	EventId    uuid.UUID `json:"eventId"`
	AttendeeId uuid.UUID `json:"attendeeId"`
}

type MeetingAttendeeUnregistered struct {
	EventId    uuid.UUID `json:"eventId"`
	AttendeeId uuid.UUID `json:"attendeeId"`
}
