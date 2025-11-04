package redis

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
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

type compressor interface {
	Compress(data []byte) ([]byte, error)
	Decompress(compressed []byte) ([]byte, error)
}

type RedisSROCache struct {
	client     *redis.Client
	serializer serializer
	compressor compressor
}

func NewRedisSROCache(addr, password string, serializer serializer, compressor compressor) (*RedisSROCache, error) {
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
		compressor: compressor,
	}, nil
}

func (c *RedisSROCache) Get(ctx context.Context, token string) (*trip.Trips, error) {
	key := generateKey(token)
	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, service.ErrNoCacheHit
		}
		return nil, err
	}

	if c.compressor != nil {
		val, err = c.compressor.Decompress(val)
		if err != nil {
			return nil, err
		}
	}

	return c.serializer.DeserializeTrips(val)
}

func (c *RedisSROCache) Set(ctx context.Context, token string, ts *trip.Trips) error {
	key := generateKey(token)
	b, err := c.serializer.SerializeTrips(ts)
	if err != nil {
		return err
	}

	if c.compressor != nil {
		b, err = c.compressor.Compress(b)
		if err != nil {
			return err
		}
	}

	return c.client.Set(ctx, key, b, ttl).Err()
}

func generateKey(token string) string {
	tokenHash := sha1.Sum([]byte(token))
	key := string(tokenHash[:])
	return fmt.Sprintf("cached_sro_%s", key)
}
