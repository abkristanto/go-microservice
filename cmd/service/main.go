package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abkristanto/go-microservice/internal/providers"
	mongorepo "github.com/abkristanto/go-microservice/internal/repositories/mongo"
	"github.com/abkristanto/go-microservice/internal/services"
	"github.com/abkristanto/go-microservice/internal/workers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LoggingProducer struct{}

func (p *LoggingProducer) Publish(ctx context.Context, payload []byte) error {
	log.Printf("Publishing message to broker: %s\n", string(payload))
	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mongoURI := envOrDefault("MONGO_URI", "mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	dbName := envOrDefault("MONGO_DB", "events_db")
	db := client.Database(dbName)

	eventCollection := envOrDefault("EVENTS_COLLECTION", "events")
	outboxCollection := envOrDefault("OUTBOX_COLLECTION", "outbox")

	eventRepo := mongorepo.NewMongoEventRepository(db, eventCollection)
	outboxRepo := mongorepo.NewMongoOutboxRepository(db, outboxCollection)
	transactionManager := mongorepo.NewMongoTransactionManager(client, db)
	apiBaseURL := envOrDefault("EVENTS_API_BASE_URL", "http://localhost:8080")
	eventProvider := providers.NewHTTPEventProvider(apiBaseURL)
	eventService := services.NewEventService(eventRepo, outboxRepo, transactionManager, eventProvider)
	eventService = services.NewLoggingService(eventService)

	producer := &LoggingProducer{}

	workers.StartSyncWorker(ctx, eventService, 30*time.Second)
	workers.StartOutboxWorker(ctx, outboxRepo, producer, 5*time.Second)

	log.Println("service started; press Ctrl+C to stop")

	<-ctx.Done()
	log.Println("shutting down")


}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
