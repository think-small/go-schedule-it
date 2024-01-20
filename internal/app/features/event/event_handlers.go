package event

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go-schedule-it/internal/app/logger"
	"log/slog"
	"net/http"
	"time"
)

func Routes(s *Service) http.Handler {
	r := chi.NewRouter()

	r.Post("/", handleEventCreated(s))

	return r
}

func handleEventCreated(s *Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var eventVM = &EventVM{}
		var event = &Event{}
		// TODO - BS - pull calendarId off URL param and verify it matches streamId on event

		err := json.NewDecoder(r.Body).Decode(&eventVM)
		if err != nil {
			slog.Error(
				"Unable to decode event object from request body.",
				slog.String(logger.INNER_ERROR, err.Error()))
			http.Error(w, "Unable to add event", http.StatusBadRequest)
			return
		}

		if eventVM.IsValid() == false {
			http.Error(w, "Invalid calendar event provided.", http.StatusBadRequest)
			return
		}

		event.Id = uuid.New()
		event.StreamId = eventVM.StreamId
		event.EventType = eventVM.EventType
		event.Timestamp = time.Now().UTC().Round(time.Microsecond)

		switch eventVM.EventType {
		case Created:
			var ec = &CalendarEventCreated{}
			err = json.Unmarshal([]byte(eventVM.Payload), &ec)
			if err != nil {
				slog.Warn(
					"Unable to decode event payload from request body.",
					slog.String(logger.INNER_ERROR, err.Error()))
				http.Error(w, "Unable to add event", http.StatusBadRequest)
				return
			}

			if ec.IsValid() == false {
				http.Error(w, "Invalid calendar event provided.", http.StatusBadRequest)
				return
			}

			marshaledJson, err := json.Marshal(ec)
			if err != nil {
				slog.Error(
					"Unable to marshal event payload into valid JSON.",
					slog.String(logger.INNER_ERROR, err.Error()))
				http.Error(w, "Unable to add calendar event.", http.StatusInternalServerError)
				return
			}
			event.Payload = marshaledJson

			//NOTE - BS - golang's time package has nanosecond precision, but postgres
			//			  timestamp with time zone data type only has microsecond precision.
			ec.ScheduledStart.UTC().Round(time.Microsecond)
			ec.ScheduledEnd.UTC().Round(time.Microsecond)

			err = s.createCalendarEvent(*event)
			if err != nil {
				http.Error(w, "Unable to insert calendar event into database.", http.StatusInternalServerError)
				return
			}

		case Canceled:
			var ec = &CalendarEventCanceled{}
			err = json.Unmarshal([]byte(eventVM.Payload), &ec)
			if err != nil {
				slog.Warn(
					"Unable to decode event payload from request body.",
					slog.String(logger.INNER_ERROR, err.Error()))
				http.Error(w, "Unable to add calendar event", http.StatusBadRequest)
				return
			}

			ec.CanceledAt = time.Now().UTC().Round(time.Microsecond)

			err = s.cancelCalendarEvent(*event)

		case Started:
			var es = &CalendarEventStarted{}
			err = json.Unmarshal([]byte(eventVM.Payload), &es)
			if err != nil {
				slog.Warn(
					"Unable to decode event payload from request body.",
					slog.String(logger.INNER_ERROR, err.Error()))
				http.Error(w, "Unable to add calendar event", http.StatusBadRequest)
				return
			}

			es.ActualStart = time.Now().UTC().Round(time.Microsecond)

			err = s.startCalendarEvent(*event)

		case Ended:
			var ee = &CalendarEventEnded{}
			err = json.Unmarshal([]byte(eventVM.Payload), &ee)
			if err != nil {
				slog.Warn(
					"Unable to decode event payload from request body.",
					slog.String(logger.INNER_ERROR, err.Error()))
				http.Error(w, "Unable to add calendar event", http.StatusBadRequest)
				return
			}

			ee.ActualEnd = time.Now().Round(time.Microsecond)

			err = s.endCalendarEvent(*event)
		}
	}
}
