package meeting

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go-schedule-it/internal/app/features/event"
	"go-schedule-it/internal/app/logger"
	"log/slog"
	"net/http"
	"time"
)

func Routes(s *Service) http.Handler {
	r := chi.NewRouter()

	r.Get("/", handleGetAllMeetings(s))
	r.Post("/", handleCreateMeeting(s))
	r.Get("/{meetingId}", handleGetMeeting(s))

	return r
}

func handleGetAllMeetings(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world"))
	}
}

func handleCreateMeeting(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var eventVM = &event.EventVM{}
		var createCalendarEvt = &event.Event{}

		slog.Debug("Deserializing event data from request body.")

		err := json.NewDecoder(r.Body).Decode(&eventVM)
		if err != nil {
			slog.Error(
				"Unable to decode event object from request body.",
				slog.String(logger.INNER_ERROR, err.Error()))
			http.Error(w, "Unable to create meeting", http.StatusBadRequest)
			return
		}

		// TODO - BS - need calendarOwnerId only for creating a meeting.
		if eventVM.IsValid() == false {
			http.Error(w, "Invalid meeting event provided.", http.StatusBadRequest)
			return
		}

		createCalendarEvt.Id = uuid.New()
		createCalendarEvt.StreamId = uuid.New()
		createCalendarEvt.EventType = eventVM.EventType
		createCalendarEvt.Timestamp = time.Now().UTC().Round(time.Microsecond)

		err = s.RegisterNewMeeting(createCalendarEvt)
		if err != nil {

		}
	}
}

func handleGetMeeting(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var meeting *Meeting

		slog.Debug("Retrieving streamId from query string.")

		rawStreamId := r.URL.Query().Get("meetingId")
		if rawStreamId == "" {
			slog.Debug("No streamId provided. Unable to hydrate meeting.")
			http.Error(w, "Unable to retrieve meeting. No streamId provided.", http.StatusBadRequest)
			return
		}
		streamId, err := uuid.Parse(rawStreamId)
		if err != nil {
			slog.Debug("Invalid streamId provided. Unable to hydrate meeting.")
			http.Error(w, "Unable to retrieve meeting. Invalid streamId.", http.StatusBadRequest)
			return
		}

		meeting, err = s.GetMeeting(streamId)
		if err != nil {
			var corruptedEventError CorruptedEventError
			switch {
			case errors.As(err, &corruptedEventError):
				http.Error(w, "Corrupted event detected. Unable to hydrate meeting.", http.StatusInternalServerError)
				return
			default:
				http.Error(w, "Unable to hydrate meeting.", http.StatusInternalServerError)
				return
			}
		}

		content, err := json.Marshal(meeting)
		if err != nil {
			slog.Error("Unable to serialize meeting.", slog.String(logger.INNER_ERROR, err.Error()))
			http.Error(w, "Unable to retrieve meeting.", http.StatusInternalServerError)
			return
		}

		_, err = w.Write(content)
		if err != nil {
			slog.Error("Unable to serialize meeting.", slog.String(logger.INNER_ERROR, err.Error()))
			http.Error(w, "Unable to serialize meeting.", http.StatusInternalServerError)
			return
		}
	}
}
