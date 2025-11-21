package models

import "time"

type Event struct {
	ID          string
	ExternalID  string
	Title       string
	Description string
	StartsAt    time.Time
}
