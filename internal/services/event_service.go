package services

import (
	"github.com/abkristanto/go-microservice/internal/models"
	"github.com/abkristanto/go-microservice/internal/repositories"
)

type EventService interface {
	SyncEvents([]models.Event) error
}

type eventService struct {
	eventRepo  repositories.EventRepository
	outboxRepo repositories.OutboxRepository
}

func NewEventService(eventRepo repositories.EventRepository, outboxRepo repositories.OutboxRepository) EventService {
	return &eventService{
		eventRepo:  eventRepo,
		outboxRepo: outboxRepo,
	}
}

func (e *eventService) SyncEvents(events []models.Event) error {
	return nil
}
