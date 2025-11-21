package workers

import (
	"context"
	"time"

	"github.com/abkristanto/go-microservice/internal/services"
)

func StartSyncWorker(ctx context.Context, svc services.EventService, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				svc.SyncEvents(ctx)
			}
		}
	}()
}
