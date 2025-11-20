package repositories

import (
	"context"

	"github.com/abkristanto/go-microservice/internal/models"
)

type OutboxRepository interface {
	Insert(ctx context.Context, outbox models.Outbox) (string, error)
	GetPending(ctx context.Context) ([]models.Event, error)
	MarkSent(ctx context.Context, id string) error
}
