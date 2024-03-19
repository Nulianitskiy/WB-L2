package event

import "time"

type Event struct {
	ID      string    `json:"event_id"`
	UserID  int64     `json:"user_id" validate:"required"`
	Data    string    `json:"data" validate:"required"`
	RawDate string    `json:"date" validate:"required"`
	Date    time.Time `json:"-"`
}

type EventRepository interface {
	Create(Event) (Event, error)
	Update(Event) (Event, error)
	Delete(int64, string) error
	GetEventsForRange(int64, time.Time, time.Time) ([]Event, error)
}
