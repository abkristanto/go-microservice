package services

import (
	"context"
	"log"
	"time"
)

type loggingService struct {
	next EventService
}

func NewLoggingService(next EventService) EventService {
	return &loggingService{
		next: next,
	}
}

func (s *loggingService) SyncEvents(ctx context.Context) error {
	start := time.Now()

	err := s.next.SyncEvents(ctx)

	took := time.Since(start)

	if err != nil {
		log.Printf("SyncEvents error: %v", err)
	} else {
		log.Printf("SyncEvents completed successfully took=%s", took)
	}

	return err
}
