package services

import (
	"context"
	"encoding/json"

	"github.com/abkristanto/go-microservice/internal/models"
	"github.com/abkristanto/go-microservice/internal/providers"
	"github.com/abkristanto/go-microservice/internal/repositories"
)

type EventService interface {
	SyncEvents(context.Context) error
}

type eventService struct {
	eventRepo     repositories.EventRepository
	outboxRepo    repositories.OutboxRepository
	tm            repositories.TransactionManager
	eventProvider providers.EventProvider
}

func NewEventService(eventRepo repositories.EventRepository, outboxRepo repositories.OutboxRepository, tm repositories.TransactionManager, eventProvider providers.EventProvider) EventService {
	return &eventService{
		eventRepo:     eventRepo,
		outboxRepo:    outboxRepo,
		tm:            tm,
		eventProvider: eventProvider,
	}
}



func (s *eventService) SyncEvents(ctx context.Context) error {
	retrievedEvents, err := s.eventProvider.GetEvents()
	if err != nil {
		return err
	}

	remoteByExternalID := make(map[string]models.Event, len(retrievedEvents))
	for _, ev := range retrievedEvents {
		remoteByExternalID[ev.ExternalID] = ev
	}

	return s.tm.WithTransaction(ctx, func(txCtx context.Context) error {
		storedEvents, err := s.eventRepo.GetEvents(txCtx)
		if err != nil {
			return err
		}

		storedByExternalID := make(map[string]models.Event, len(storedEvents))
		for _, stored := range storedEvents {
			storedByExternalID[stored.ExternalID] = stored
		}

		for _, remote := range retrievedEvents {
			stored, exists := storedByExternalID[remote.ExternalID]

			if !exists {
				if err := s.syncEventAndEnqueue(txCtx, remote, "event.created"); err != nil {
					return err
				}
				continue
			}

			if hasChanged(stored, remote) {
				if err := s.syncEventAndEnqueue(txCtx, remote, "event.updated"); err != nil {
					return err
				}
			}
		}

		for _, stored := range storedEvents {
			_, exists := remoteByExternalID[stored.ExternalID]
			if exists {
				continue
			}

			if err := s.eventRepo.DeleteEventByID(txCtx, stored.ID); err != nil {
				return err
			}

			p := payload{
				ChangeType:       "event.deleted",
				APISource:        s.eventProvider.APISource(),
				ResourceLocation: "",
				Event:            stored,
			}

			payloadBytes, err := json.Marshal(p)
			if err != nil {
				return err
			}

			outbox := models.Outbox{
				Status:     "pending",
				Payload:    payloadBytes,
				RetryCount: 0,
			}

			if err := s.outboxRepo.Insert(txCtx, outbox); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *eventService) syncEventAndEnqueue(ctx context.Context, remote models.Event, changeType string) error {
	location, err := s.eventRepo.UpsertEvent(ctx, remote)
	if err != nil {
		return err
	}

	p := payload{
		ChangeType:       changeType,
		APISource:        s.eventProvider.APISource(),
		ResourceLocation: location,
		Event:            remote,
	}

	payloadBytes, err := json.Marshal(p)
	if err != nil {
		return err
	}

	outbox := models.Outbox{
		Status:     "pending",
		Payload:    payloadBytes,
		RetryCount: 0,
	}

	return s.outboxRepo.Insert(ctx, outbox)
}
