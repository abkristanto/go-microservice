package mongo

import (
	"context"

	"github.com/abkristanto/go-microservice/internal/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoTransactionManager struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoTransactionManager(client *mongo.Client, db *mongo.Database) repositories.TransactionManager {
	return &MongoTransactionManager{
		client: client,
		db:     db,
	}
}

func (tm *MongoTransactionManager) WithTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {

	session, err := tm.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	}

	_, err = session.WithTransaction(
		ctx,
		callback,
		options.Transaction().
			SetReadConcern(tm.db.ReadConcern()).
			SetWriteConcern(tm.db.WriteConcern()),
	)

	return err
}
