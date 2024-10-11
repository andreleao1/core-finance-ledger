package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache() *RedisCache {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &RedisCache{
		client: client,
		ctx:    ctx,
	}
}

func (r *RedisCache) SaveBitcoinPrice(price float64) error {
	key := "bitcoin:price"
	value := fmt.Sprintf("%.2f", price)

	err := r.client.Set(r.ctx, key, value, 1*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to save bitcoin price to redis: %w", err)
	}
	return nil
}

func (r *RedisCache) GetBitcoinPrice() (float64, error) {
	key := "bitcoin:price"

	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return 0, fmt.Errorf("no bitcoin price found in redis")
	} else if err != nil {
		return 0, fmt.Errorf("failed to get bitcoin price from redis: %w", err)
	}

	price, err := strconv.ParseFloat(val, 64)

	if err != nil {
		return 0, fmt.Errorf("failed to parse bitcoin price from redis: %w", err)
	}

	return price, nil
}
