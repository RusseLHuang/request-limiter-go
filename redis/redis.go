package redis

import "github.com/go-redis/redis/v8"

func NewClient(
	url string,
	password string,
) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
	})

	return client
}
