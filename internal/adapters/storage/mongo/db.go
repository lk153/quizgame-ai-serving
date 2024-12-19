package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	"github.com/lk153/quizgame-ai-serving/internal/adapters/config"
)

type DB struct {
	*mongo.Client
	DB  *mongo.Database
	url string
}

func New(ctx context.Context, config *config.DB) (*DB, error) {
	url := fmt.Sprintf("%s://%s:%s@%s/?retryWrites=true&w=majority&appName=%s",
		config.Connection,
		config.User,
		config.Password,
		config.Host,
		config.ClusterName,
	)

	bsonOpts := &options.BSONOptions{
		UseJSONStructTags: true,
		NilSliceAsEmpty:   true,
		NilMapAsEmpty:     true,
	}
	client, err := mongo.Connect(options.Client().ApplyURI(url).SetBSONOptions(bsonOpts))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	db := client.Database(config.DBName)

	return &DB{
		client,
		db,
		url,
	}, nil
}
