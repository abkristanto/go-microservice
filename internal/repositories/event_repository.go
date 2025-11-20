package repositories

import (
	"context"

	"github.com/abkristanto/go-microservice/internal/models"
)

type EventRepository interface {
	GetEvents(ctx context.Context) ([]models.Event, error)
	UpsertEvent(ctx context.Context, event models.Event) (string, error) 
}
