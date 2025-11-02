package redis

import (
	"context"
	"crypto/sha1"
	"errors"
	"time"

	"github.com/de4et/flight-booking/internal/model/trip"
	"github.com/de4et/flight-booking/internal/service"

	"github.com/redis/go-redis/v9"
)

const (
	ttl     = time.Minute * 15
	timeout = time.Second
)

type serializer interface {
	SerializeTrips(ts *trip.Trips) ([]byte, error)
	DeserializeTrips(data []byte) (*trip.Trips, error)
}

type RedisSROCache struct {
	client     *redis.Client
	serializer serializer
}

func NewRedisSROCache(addr, password string, serializer serializer) (*RedisSROCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisSROCache{
		client:     rdb,
		serializer: serializer,
	}, nil
}

func (c *RedisSROCache) Get(ctx context.Context, token string) (*trip.Trips, error) {
	tokenHash := sha1.Sum([]byte(token))
	key := string(tokenHash[:])

	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, service.ErrNoCacheHit
		}
		return nil, err
	}

	return c.serializer.DeserializeTrips(val)
}

func (c *RedisSROCache) Set(ctx context.Context, token string, ts *trip.Trips) error {
	tokenHash := sha1.Sum([]byte(token))
	key := string(tokenHash[:])

	b, err := c.serializer.SerializeTrips(ts)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, b, 0).Err()
}
