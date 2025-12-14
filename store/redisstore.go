package store

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
    client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
    return &RedisStore{client: client}
}

func (r *RedisStore) Get(ctx context.Context, key string) (string, error) {
    return r.client.Get(ctx, key).Result()
}

func (r *RedisStore) Set(ctx context.Context, key string, value interface{}) error {
    var data string

    switch v := value.(type) {
    case string:
        data = v
    default:
        b, err := json.Marshal(v)
        if err != nil {
            return err
        }
        data = string(b)
    }

    return r.client.Set(ctx, key, data, 0).Err()
}
