package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go-schedule-it/internal/app/db"
	"go-schedule-it/internal/app/features/calendar"
	"go-schedule-it/internal/app/features/event"
	"go-schedule-it/internal/app/logger"
	"log/slog"
	"net/http"
	"os"
)

type server struct {
}

type serverConfig struct {
	port         string
	dbProvider   string
	dbConnString string
}

func Run() {

	cfg := &serverConfig{
		port:         os.Getenv("APP_PORT"),
		dbProvider:   os.Getenv("DB_PROVIDER"),
		dbConnString: os.Getenv("DB_CONNECTION_STRING"),
	}

	server := &server{}
	fmt.Println(server)

	eventStreamWriter := db.NewEventStreamWriter(cfg.dbProvider, cfg.dbConnString)
	eventStreamReader := db.NewEventStreamReader(cfg.dbProvider, cfg.dbConnString)
	eventService := event.NewEventService(eventStreamWriter, eventStreamReader)
	calendarService := calendar.NewCalendarService(eventStreamReader)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.Heartbeat("/health"))

	router.Mount("/calendars", calendar.Routes(calendarService))
	router.Mount("/calendars/{calendarId}/events", event.Routes(eventService))

	slog.Info(fmt.Sprintf("Starting http server on port: %s\n", cfg.port))

	err := http.ListenAndServe(cfg.port, router)
	if err != nil {
		slog.Error("Unable to start http server.", slog.String(logger.INNER_ERROR, err.Error()))
		os.Exit(1)
	}
}
