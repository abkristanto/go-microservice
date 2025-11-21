package documents

import "time"

type EventDocument struct {
	ID          string    `bson:"_id"`
	ExternalID  string    `bson:"external_id"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	StartsAt    time.Time `bson:"starts_at"`
}
