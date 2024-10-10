package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache khởi tạo kết nối Redis
func NewRedisCache(host string, port int, password string, db int) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		client: rdb,
		ctx:    context.Background(),
	}
}

// Set lưu trữ dữ liệu vào Redis với thời gian hết hạn
func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

// Get lấy dữ liệu từ Redis bằng key
func (r *RedisCache) Get(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key does not exist")
	}
	return val, err
}

// Delete xóa một key khỏi Redis
func (r *RedisCache) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}
