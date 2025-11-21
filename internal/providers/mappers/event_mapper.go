package mappers

import (
	"github.com/abkristanto/go-microservice/internal/models"
	"github.com/abkristanto/go-microservice/internal/providers/dtos"
)

func ToDomainEvent(dto dtos.Event) models.Event {
	return models.Event{
		ExternalID:  dto.ExternalID,
		Title:       dto.Title,
		Description: dto.Description,
		StartsAt:    dto.StartsAt,
	}
}
