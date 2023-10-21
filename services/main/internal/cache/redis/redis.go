package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/raisultan/url-shortener/services/main/internal/config"
	"github.com/raisultan/url-shortener/services/main/internal/lib/logger/sl"
	"golang.org/x/exp/slog"
	"time"
)

const urlTTL = 24 * time.Hour

type Cache struct {
	client *redis.Client
}

func New(config config.Cache, ctx context.Context) (*Cache, error) {
	const op = "cache.redis.New"

	options, err := redis.ParseURL(config.URL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	client := redis.NewClient(options)
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Cache{client}, nil
}

func (c *Cache) Close(log *slog.Logger) {
	err := c.client.Close()
	if err != nil {
		log.Error("could not close cache", sl.Err(err))
	}
}

func (c *Cache) SaveUrl(ctx context.Context, urlToSave string, alias string) error {
	const op = "cache.redis.SaveUrl"

	err := c.client.Set(ctx, alias, urlToSave, urlTTL).Err()
	if err != nil {
		return fmt.Errorf("%s: could not save url to cache %w", op, err)
	}

	return nil
}

func (c *Cache) GetUrl(ctx context.Context, alias string) (string, error) {
	const op = "cache.redis.GetUrl"

	url, err := c.client.Get(ctx, alias).Result()
	if err != nil {
		return "", fmt.Errorf("%s: could not get url from cache %w", op, err)
	}

	return url, nil
}

func (c *Cache) DeleteUrl(ctx context.Context, alias string) error {
	const op = "cache.redis.DeleteUrl"

	err := c.client.Del(ctx, alias).Err()
	if err != nil {
		return fmt.Errorf("%s: could not delete url from cache %w", op, err)
	}

	return nil
}
