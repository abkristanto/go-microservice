package models

import "time"

type Event struct {
	ID          string
	Title       string
	Description string
	StartsAt    time.Time
}
