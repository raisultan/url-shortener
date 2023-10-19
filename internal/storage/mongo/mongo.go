package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/raisultan/url-shortener/internal/lib/logger/sl"
	"github.com/raisultan/url-shortener/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slog"
)

type Storage struct {
	db *mongo.Collection
}

func New(
	uri string,
	databaseName string,
	collectionName string,
	ctx context.Context,
) (*Storage, error) {
	const op = "storage.mongo.New"

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: client.Database(databaseName).Collection(collectionName),
	}, nil
}

func (s *Storage) Close(ctx context.Context, log *slog.Logger) {
	err := s.db.Database().Client().Disconnect(ctx)
	if err != nil {
		log.Error("could not close storage", sl.Err(err))
	}
}

func (s *Storage) SaveUrl(ctx context.Context, urlToSave, alias string) error {
	record := struct {
		Alias string `bson:"alias"`
		Url   string `bson:"url"`
	}{
		Alias: alias,
		Url:   urlToSave,
	}

	_, err := s.db.InsertOne(ctx, record)
	if err != nil {
		return fmt.Errorf("failed to save url with the alias %s: %w", alias, err)
	}

	return nil
}

func (s *Storage) GetUrl(ctx context.Context, alias string) (string, error) {
	var result struct {
		Url string `bson:"url"`
	}

	err := s.db.FindOne(ctx, bson.D{{"alias", alias}}).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", storage.ErrUrlNotFound
	}
	if err != nil {
		return "", fmt.Errorf("failed to find document with the alias %s: %w", alias, err)
	}

	return result.Url, nil
}
