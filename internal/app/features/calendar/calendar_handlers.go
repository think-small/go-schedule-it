package calendar

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go-schedule-it/internal/app/logger"
	"log/slog"
	"net/http"
)

func Routes(s *Service) http.Handler {
	r := chi.NewRouter()

	r.Get("/", handleGetAllCalendars(s))
	r.Get("/{calendarId}", handleGetCalendar(s))

	return r
}

func handleGetAllCalendars(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world"))
	}
}

func handleGetCalendar(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var calendar *Calendar

		slog.Debug("Retrieving streamId from query string.")

		rawStreamId := r.URL.Query().Get("calendarId")
		if rawStreamId == "" {
			slog.Debug("No streamId provided. Unable to hydrate calendar.")
			http.Error(w, "Unable to retrieve calendar. No streamId provided.", http.StatusBadRequest)
			return
		}
		streamId, err := uuid.Parse(rawStreamId)
		if err != nil {
			slog.Debug("Invalid streamId provided. Unable to hydrate calendar.")
			http.Error(w, "Unable to retrieve calendar. Invalid streamId.", http.StatusBadRequest)
			return
		}

		calendar, err = s.GetCalendar(streamId)
		if err != nil {
			var corruptedEventError CorruptedEventError
			switch {
			case errors.As(err, &corruptedEventError):
				http.Error(w, "Corrupted event detected. Unable to hydrate calendar.", http.StatusInternalServerError)
				return
			default:
				http.Error(w, "Unable to hydrate calendar.", http.StatusInternalServerError)
				return
			}
		}

		content, err := json.Marshal(calendar)
		if err != nil {
			slog.Error("Unable to serialize calendar.", slog.String(logger.INNER_ERROR, err.Error()))
			http.Error(w, "Unable to retrieve calendar.", http.StatusInternalServerError)
			return
		}

		_, err = w.Write([]byte(content))
		if err != nil {
			slog.Error("Unable to serialize calendar.", slog.String(logger.INNER_ERROR, err.Error()))
			http.Error(w, "Unable to serialize calendar.", http.StatusInternalServerError)
			return
		}
	}
}
