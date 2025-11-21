package mocks

import (
	"context"

	"github.com/abkristanto/go-microservice/internal/models"
	"github.com/abkristanto/go-microservice/internal/providers"
	"github.com/abkristanto/go-microservice/internal/repositories"
	"github.com/abkristanto/go-microservice/internal/services"
)

type EventProviderMock struct {
	Events []models.Event
	Err    error
	Source string

	Calls int
}

var _ providers.EventProvider = (*EventProviderMock)(nil)

func (m *EventProviderMock) GetEvents() ([]models.Event, error) {
	m.Calls++
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Events, nil
}

func (m *EventProviderMock) APISource() string {
	if m.Source != "" {
		return m.Source
	}
	return "mock-event-provider"
}

type EventRepositoryMock struct {
	StoredEvents []models.Event

	UpsertedEvents []models.Event
	DeletedIDs     []string

	GetEventsErr   error
	UpsertEventErr error
	DeleteErr      error
}

var _ repositories.EventRepository = (*EventRepositoryMock)(nil)

func (m *EventRepositoryMock) GetEvents(ctx context.Context) ([]models.Event, error) {
	if m.GetEventsErr != nil {
		return nil, m.GetEventsErr
	}
	return m.StoredEvents, nil
}

func (m *EventRepositoryMock) UpsertEvent(ctx context.Context, ev models.Event) (string, error) {
	if m.UpsertEventErr != nil {
		return "", m.UpsertEventErr
	}
	m.UpsertedEvents = append(m.UpsertedEvents, ev)

	// simulate a location string
	if ev.ID == "" {
		return "mockdb.events.mock-id", nil
	}
	return "mockdb.events." + ev.ID, nil
}

func (m *EventRepositoryMock) DeleteEventByID(ctx context.Context, id string) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	m.DeletedIDs = append(m.DeletedIDs, id)
	return nil
}

type OutboxRepositoryMock struct {
	Pending   []models.Outbox
	Inserted  []models.Outbox
	MarkedIDs []string

	GetPendingErr error
	InsertErr     error
	MarkSentErr   error
}

var _ repositories.OutboxRepository = (*OutboxRepositoryMock)(nil)

func (m *OutboxRepositoryMock) GetPending(ctx context.Context) ([]models.Outbox, error) {
	if m.GetPendingErr != nil {
		return nil, m.GetPendingErr
	}
	return m.Pending, nil
}

func (m *OutboxRepositoryMock) Insert(ctx context.Context, o models.Outbox) error {
	if m.InsertErr != nil {
		return m.InsertErr
	}
	m.Inserted = append(m.Inserted, o)
	return nil
}

func (m *OutboxRepositoryMock) MarkSent(ctx context.Context, id string) error {
	if m.MarkSentErr != nil {
		return m.MarkSentErr
	}
	m.MarkedIDs = append(m.MarkedIDs, id)
	return nil
}

type TransactionManagerMock struct {
	Calls int
	Err   error
}

var _ repositories.TransactionManager = (*TransactionManagerMock)(nil)

func (m *TransactionManagerMock) WithTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	m.Calls++
	if m.Err != nil {
		return m.Err
	}
	return fn(ctx)
}

type ProducerMock struct {
	Payloads [][]byte
	Err      error
}

var _ workersProducer = (*ProducerMock)(nil)

type workersProducer interface {
	Publish(ctx context.Context, payload []byte) error
}

func (m *ProducerMock) Publish(ctx context.Context, payload []byte) error {
	if m.Err != nil {
		return m.Err
	}
	m.Payloads = append(m.Payloads, payload)
	return nil
}

type EventServiceMock struct {
	Calls int
	Err   error
}

var _ services.EventService = (*EventServiceMock)(nil)

func (m *EventServiceMock) SyncEvents(ctx context.Context) error {
	m.Calls++
	return m.Err
}
