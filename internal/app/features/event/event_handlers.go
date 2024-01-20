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
		var eventVM = &ViewModel{}
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
			http.Error(w, "Invalid meeting event provided.", http.StatusBadRequest)
			return
		}

		event.Id = uuid.New()
		event.StreamId = eventVM.StreamId
		event.EventType = eventVM.EventType
		event.Timestamp = time.Now().UTC().Round(time.Microsecond)

		switch eventVM.EventType {
		case Created:
			var meetingCreated = &MeetingCreated{}
			err = json.Unmarshal([]byte(eventVM.Payload), &meetingCreated)
			if err != nil {
				slog.Warn(
					"Unable to decode event payload from request body.",
					slog.String(logger.INNER_ERROR, err.Error()))
				http.Error(w, "Unable to add event", http.StatusBadRequest)
				return
			}

			if meetingCreated.IsValid() == false {
				http.Error(w, "Invalid meeting event provided.", http.StatusBadRequest)
				return
			}

			marshaledJson, err := json.Marshal(meetingCreated)
			if err != nil {
				slog.Error(
					"Unable to marshal event payload into valid JSON.",
					slog.String(logger.INNER_ERROR, err.Error()))
				http.Error(w, "Unable to add meeting event.", http.StatusInternalServerError)
				return
			}
			event.Payload = marshaledJson

			//NOTE - BS - golang's time package has nanosecond precision, but postgres
			//			  timestamp with time zone data type only has microsecond precision.
			meetingCreated.ScheduledStart.UTC().Round(time.Microsecond)
			meetingCreated.ScheduledEnd.UTC().Round(time.Microsecond)

			err = s.createMeeting(*event)
			if err != nil {
				http.Error(w, "Unable to insert meeting event into database.", http.StatusInternalServerError)
				return
			}

		case Canceled:
			var meetingCanceled = &MeetingCanceled{}
			err = json.Unmarshal([]byte(eventVM.Payload), &meetingCanceled)
			if err != nil {
				slog.Warn(
					"Unable to decode event payload from request body.",
					slog.String(logger.INNER_ERROR, err.Error()))
				http.Error(w, "Unable to add meeting event", http.StatusBadRequest)
				return
			}

			meetingCanceled.CanceledAt = time.Now().UTC().Round(time.Microsecond)

			err = s.cancelMeeting(*event)

		case Started:
			var meetingStarted = &MeetingStarted{}
			err = json.Unmarshal([]byte(eventVM.Payload), &meetingStarted)
			if err != nil {
				slog.Warn(
					"Unable to decode event payload from request body.",
					slog.String(logger.INNER_ERROR, err.Error()))
				http.Error(w, "Unable to add meeting event", http.StatusBadRequest)
				return
			}

			meetingStarted.ActualStart = time.Now().UTC().Round(time.Microsecond)

			err = s.startMeeting(*event)

		case Ended:
			var meetingEnded = &MeetingEnded{}
			err = json.Unmarshal([]byte(eventVM.Payload), &meetingEnded)
			if err != nil {
				slog.Warn(
					"Unable to decode event payload from request body.",
					slog.String(logger.INNER_ERROR, err.Error()))
				http.Error(w, "Unable to add meeting event", http.StatusBadRequest)
				return
			}

			meetingEnded.ActualEnd = time.Now().Round(time.Microsecond)

			err = s.endMeeting(*event)
		}
	}
}
