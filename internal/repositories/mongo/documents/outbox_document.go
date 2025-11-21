package documents

import "time"

type OutboxDocument struct {
	ID          string     `bson:"_id"`
	Status      string     `bson:"status"`
	Payload     []byte     `bson:"payload"`
	RetryCount  int        `bson:"retry_count"`
	CreatedAt   time.Time  `bson:"created_at"`
	ProcessedAt *time.Time `bson:"processed_at,omitempty"`
}
