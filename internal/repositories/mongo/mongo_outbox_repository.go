package mongo

import (
	"context"
	"time"

	"github.com/abkristanto/go-microservice/internal/models"
	"github.com/abkristanto/go-microservice/internal/repositories"
	"github.com/abkristanto/go-microservice/internal/repositories/mongo/documents"
	"github.com/abkristanto/go-microservice/internal/repositories/mongo/mappers"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoOutboxRepository struct {
	coll *mongo.Collection
}

func NewMongoOutboxRepository(db *mongo.Database, collection string) repositories.OutboxRepository {
	return &MongoOutboxRepository{
		coll: db.Collection(collection),
	}
}

func (r *MongoOutboxRepository) Insert(ctx context.Context, outbox models.Outbox) error {

	if outbox.ID == "" {
		outbox.ID = uuid.NewString()
	}

	if outbox.CreatedAt.IsZero() {
		outbox.CreatedAt = time.Now()
	}
	
	doc := mappers.ToOutboxDocument(outbox)

	_, err := r.coll.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoOutboxRepository) GetPending(ctx context.Context) ([]models.Outbox, error) {
	filter := bson.M{"status": "pending"}

	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []documents.OutboxDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	out := make([]models.Outbox, len(docs))
	for i, d := range docs {
		out[i] = mappers.ToDomainOutbox(d)
	}

	return out, nil
}

func (r *MongoOutboxRepository) MarkSent(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status":       "sent",
			"processed_at": time.Now(),
		},
	}

	_, err := r.coll.UpdateOne(ctx, filter, update)
	return err
}
