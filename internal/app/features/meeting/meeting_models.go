package meeting

import (
	"encoding/json"
	"github.com/google/uuid"
	"go-schedule-it/internal/app/features/event"
	"go-schedule-it/internal/app/logger"
	"log/slog"
	"time"
)

type CorruptedEventError struct {
	event event.Event
}

func (ce CorruptedEventError) Error() string {
	return "Corrupted event found; unable to hydrate meeting."
}

type Meeting struct {
	Id             *uuid.UUID
	CalendarId     *uuid.UUID
	CreatedAt      *time.Time
	ScheduledStart *time.Time
	ScheduledEnd   *time.Time
	ActualStart    *time.Time
	ActualEnd      *time.Time
	CanceledAt     *time.Time
	HostId         *uuid.UUID
	Attendees      []uuid.UUID
	events         []event.Event
}

func (c *Meeting) Apply(events []event.Event) error {
	for _, e := range events {
		switch e.EventType {
		case event.Created:
			createdEvt := &event.MeetingCreated{}
			err := json.Unmarshal(e.Payload, &createdEvt)
			if err != nil {
				corruptedEventErr := &CorruptedEventError{
					event: e,
				}
				slog.Error(corruptedEventErr.Error(), slog.String(logger.INNER_ERROR, err.Error()))
				return corruptedEventErr
			}

			if createdEvt.IsValid() == false {
				corruptedEventErr := &CorruptedEventError{
					event: e,
				}
				slog.Error(corruptedEventErr.Error())
				return corruptedEventErr
			}

			c.CreatedAt = &e.Timestamp
			c.ScheduledStart = &createdEvt.ScheduledStart
			c.ScheduledEnd = &createdEvt.ScheduledEnd
			c.CalendarId = &createdEvt.CalendarId
		case event.Canceled:
			canceledEvt := &event.MeetingCanceled{}
			err := json.Unmarshal(e.Payload, &canceledEvt)
			if err != nil {
				corruptedEventErr := &CorruptedEventError{event: e}
				slog.Error(corruptedEventErr.Error(), slog.String(logger.INNER_ERROR, err.Error()))
				return corruptedEventErr
			}

			c.CanceledAt = &canceledEvt.CanceledAt
		case event.Started:
			startedEvt := &event.MeetingStarted{}
			err := json.Unmarshal(e.Payload, &startedEvt)
			if err != nil {
				corruptedEventErr := CorruptedEventError{event: e}
				slog.Error(corruptedEventErr.Error(), slog.String(logger.INNER_ERROR, err.Error()))
				return corruptedEventErr
			}

			c.ActualStart = &startedEvt.ActualStart
		case event.Ended:
			endedEvt := &event.MeetingEnded{}
			err := json.Unmarshal(e.Payload, &endedEvt)
			if err != nil {
				corruptedEventErr := CorruptedEventError{event: e}
				slog.Error(corruptedEventErr.Error(), slog.String(logger.INNER_ERROR, err.Error()))
				return corruptedEventErr
			}

			c.ActualEnd = &endedEvt.ActualEnd
		case event.AttendeeRegistered:
			registeredEvt := &event.MeetingAttendeeRegistered{}
			err := json.Unmarshal(e.Payload, registeredEvt)
			if err != nil {
				corruptedEventErr := CorruptedEventError{event: e}
				slog.Error(corruptedEventErr.Error(), slog.String(logger.INNER_ERROR, err.Error()))
				return corruptedEventErr
			}

			var isAlreadyRegistered = false
			for _, a := range c.Attendees {
				if a == registeredEvt.AttendeeId {
					isAlreadyRegistered = true
				}
			}

			if isAlreadyRegistered {
				corruptedEventErr := &CorruptedEventError{event: e}
				slog.Error(corruptedEventErr.Error(), slog.String(logger.INNER_ERROR, "Duplicate attendee."))
				return corruptedEventErr
			}

			c.Attendees = append(c.Attendees, registeredEvt.AttendeeId)
		case event.AttendeeUnregistered:
			unregisteredEvt := &event.MeetingAttendeeUnregistered{}
			err := json.Unmarshal(e.Payload, &unregisteredEvt)
			if err != nil {
				corruptedEventErr := &CorruptedEventError{event: e}
				slog.Error(corruptedEventErr.Error(), slog.String(logger.INNER_ERROR, err.Error()))
				return corruptedEventErr
			}

			currIndex := 0
			newAttendees := make([]uuid.UUID, len(c.Attendees)-1)
			for _, a := range c.Attendees {
				if a != unregisteredEvt.AttendeeId {
					newAttendees[currIndex] = a
					currIndex++
				}
			}

			c.Attendees = newAttendees
		case event.Unsupported:
			corruptedEventErr := &CorruptedEventError{event: e}
			slog.Error(corruptedEventErr.Error(), slog.String(logger.INNER_ERROR, "Unsupported event type"))
			return corruptedEventErr
		}
	}
	return nil
}
