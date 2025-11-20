package models

import "time"

type Outbox struct {
	ID          string
	Status      string
	Payload     []byte
	RetryCount  int
	CreatedAt   time.Time
	ProcessedAt time.Time
}
