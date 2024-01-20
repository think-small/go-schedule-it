package calendar

import (
	"github.com/google/uuid"
	"go-schedule-it/internal/app/features/meeting"
	"time"
)

type Calendar struct {
	id             uuid.UUID
	ownerId        uuid.UUID
	allowedUserIds []uuid.UUID
	meetings       []meeting.Meeting
	createdAt      time.Time
	archivedAt     *time.Time
}
