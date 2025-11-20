package dtos

import "time"

type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartsAt    time.Time `json:"starts_at"`
}
