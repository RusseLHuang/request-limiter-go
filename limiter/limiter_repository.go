package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type LimiterRepository interface {
	Increment(ctx context.Context, key string) int
	Get(ctx context.Context, key string) string
	Set(ctx context.Context, key string, value string, duration int)
	Del(ctx context.Context, key string)
}

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(redisClient *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: redisClient,
	}
}

func (redisRepository RedisRepository) Increment(
	ctx context.Context,
	key string,
) int {
	res, err := redisRepository.client.Incr(ctx, key).Result()
	if err != nil {
		panic(err)
	}

	return int(res)
}

func (redisRepository RedisRepository) Get(
	ctx context.Context,
	key string,
) string {
	res, err := redisRepository.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return ""
	}

	return res
}

func (redisRepository RedisRepository) Set(
	ctx context.Context,
	key string,
	value string,
	duration int,
) {
	expirationDuration := time.Duration(duration) * time.Duration(time.Second)

	res := redisRepository.client.Set(
		ctx,
		key,
		value,
		expirationDuration,
	)

	if res.Err() != nil {
		panic(res.Err())
	}
}

func (redisRepository RedisRepository) Del(
	ctx context.Context,
	key string,
) {
	res := redisRepository.client.Del(ctx, key)

	if res.Err() != nil {
		panic(res.Err())
	}
}
