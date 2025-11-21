package mongo

import (
	"context"

	"github.com/abkristanto/go-microservice/internal/models"
	"github.com/abkristanto/go-microservice/internal/repositories"
	"github.com/abkristanto/go-microservice/internal/repositories/mongo/documents"
	"github.com/abkristanto/go-microservice/internal/repositories/mongo/mappers"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoEventRepository struct {
	coll *mongo.Collection
}

func NewMongoEventRepository(db *mongo.Database, collection string) repositories.EventRepository {
	return &MongoEventRepository{
		coll: db.Collection(collection),
	}
}

func (r *MongoEventRepository) GetEvents(ctx context.Context) ([]models.Event, error) {
	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []documents.EventDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	events := make([]models.Event, len(docs))
	for i, d := range docs {
		events[i] = mappers.ToDomainEvent(d)
	}

	return events, nil
}

func (r *MongoEventRepository) UpsertEvent(ctx context.Context, event models.Event) (string, error) {
	if event.ID == "" {
		event.ID = uuid.NewString()
	}
	doc := mappers.ToEventDocument(event)

	filter := bson.M{"_id": doc.ID}
	update := bson.M{"$set": doc}

	opts := options.Update().SetUpsert(true)

	_, err := r.coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return "", err
	}

	location := r.coll.Database().Name() + "." + r.coll.Name() + "." + doc.ID

	return location, nil
}

func (r *MongoEventRepository) DeleteEventByID(ctx context.Context, id string) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
