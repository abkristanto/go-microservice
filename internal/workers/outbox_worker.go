// internal/workers/outbox_worker.go
package workers

import (
	"context"
	"log"
	"time"

	"github.com/abkristanto/go-microservice/internal/repositories"
)

type Producer interface {
	Publish(ctx context.Context, payload []byte) error
}

func StartOutboxWorker(ctx context.Context, outboxRepo repositories.OutboxRepository, producer Producer, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := processPending(ctx, outboxRepo, producer); err != nil {
					log.Printf("outbox worker error: %v", err)
				}
			}
		}
	}()
}

func processPending(ctx context.Context, outboxRepo repositories.OutboxRepository, producer Producer) error {
	entries, err := outboxRepo.GetPending(ctx)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if err := producer.Publish(ctx, entry.Payload); err != nil {
			// increment RetryCount, log, etc.
			continue
		}

		if err := outboxRepo.MarkSent(ctx, entry.ID); err != nil {
			return err
		}
	}

	return nil
}
