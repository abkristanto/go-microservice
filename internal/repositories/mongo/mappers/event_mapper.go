package mappers

import (
	"github.com/abkristanto/go-microservice/internal/models"
	"github.com/abkristanto/go-microservice/internal/repositories/mongo/documents"
)

func ToEventDocument(e models.Event) documents.EventDocument {
	return documents.EventDocument{
		ID:          e.ID,
		ExternalID:  e.ExternalID,
		Title:       e.Title,
		Description: e.Description,
		StartsAt:    e.StartsAt,
	}
}

func ToDomainEvent(d documents.EventDocument) models.Event {
	return models.Event{
		ID:          d.ID,
		ExternalID:  d.ExternalID,
		Title:       d.Title,
		Description: d.Description,
		StartsAt:    d.StartsAt,
	}
}
