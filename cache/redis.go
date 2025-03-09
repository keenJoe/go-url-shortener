package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/keenJoe/go-url-shortener/config"
)

var RedisClient *redis.Client
var ctx = context.Background()

// InitRedis 初始化Redis连接
func InitRedis(config *config.Config) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
		PoolSize: 100,
	})

	_, err := RedisClient.Ping(ctx).Result()
	return err
}

// SetURL 缓存URL映射
func SetURL(shortCode, originalURL string, expiration time.Duration) error {
	return RedisClient.Set(ctx, "url:"+shortCode, originalURL, expiration).Err()
}

// GetURL 获取URL映射
func GetURL(shortCode string) (string, error) {
	return RedisClient.Get(ctx, "url:"+shortCode).Result()
}

// DeleteURL 删除URL映射
func DeleteURL(shortCode string) error {
	return RedisClient.Del(ctx, "url:"+shortCode).Err()
}

// IncrementCounter 增加访问计数
func IncrementCounter(shortCode string) error {
	return RedisClient.Incr(ctx, "counter:"+shortCode).Err()
}

// GetCounter 获取访问计数
func GetCounter(shortCode string) (int64, error) {
	return RedisClient.Get(ctx, "counter:"+shortCode).Int64()
}
