package mappers

import (
	"github.com/abkristanto/go-microservice/internal/models"
	"github.com/abkristanto/go-microservice/internal/repositories/mongo/documents"
)

func ToOutboxDocument(o models.Outbox) documents.OutboxDocument {
	return documents.OutboxDocument{
		ID:          o.ID,
		Status:      o.Status,
		Payload:     o.Payload,
		RetryCount:  o.RetryCount,
		CreatedAt:   o.CreatedAt,
		ProcessedAt: o.ProcessedAt,
	}
}

func ToDomainOutbox(d documents.OutboxDocument) models.Outbox {
	return models.Outbox{
		ID:          d.ID,
		Status:      d.Status,
		Payload:     d.Payload,
		RetryCount:  d.RetryCount,
		CreatedAt:   d.CreatedAt,
		ProcessedAt: d.ProcessedAt,
	}
}
