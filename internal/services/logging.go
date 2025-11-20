package services

import (
	"log"
	"time"

	"github.com/abkristanto/go-microservice/internal/models"
)

type loggingService struct {
	next EventService
}

func NewLoggingService(next EventService) EventService {
	return &loggingService{
		next: next,
	}
}

func (s *loggingService) SyncEvents(events []models.Event) error {
	start := time.Now()
	log.Printf("SyncEvents called with %d events", len(events))

	err := s.next.SyncEvents(events)

	took := time.Since(start)

	if err != nil {
		log.Printf("SyncEvents error: %v", err)
	} else {
		log.Printf("SyncEvents completed successfully took=%s", took)
	}

	return err
}
