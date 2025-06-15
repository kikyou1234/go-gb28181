package redis

import (
	"context"

	"github.com/redis/go-redis/v9"

	"time"
)

const (
	REDIS_POOL_SIZE         = 100
	REDIS_MIN_IDLE_CONNS    = 10
	REDIS_POOL_TIMEOUT_SEC  = 5
	REDIS_DIAL_TIMEOUT_SEC  = 5
	REDIS_READ_TIMEOUT_SEC  = 3
	REDIS_WRITE_TIMEOUT_SEC = 3
)

var (
	RDB *redis.Client
)

func InitRedis(addr, password string, db int) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,

		PoolSize:     REDIS_POOL_SIZE,
		MinIdleConns: REDIS_MIN_IDLE_CONNS,
		PoolTimeout:  time.Duration(REDIS_POOL_TIMEOUT_SEC) * time.Second,
		DialTimeout:  time.Duration(REDIS_DIAL_TIMEOUT_SEC) * time.Second,
		ReadTimeout:  time.Duration(REDIS_READ_TIMEOUT_SEC) * time.Second,
		WriteTimeout: time.Duration(REDIS_WRITE_TIMEOUT_SEC) * time.Second,
	})

	// 测试连接
	ctx := context.Background()
	if _, err := RDB.Ping(ctx).Result(); err != nil {
		panic(err)
	}

}

func GetRedisClient() *redis.Client {
	if RDB == nil {
		panic("current redis client is not set")
	}
	return RDB
}

func HSet(ctx context.Context, redis_key string, key, value string) error {
	rdb := GetRedisClient()
	return rdb.HSet(ctx, redis_key, key, value).Err()
}

func HDel(ctx context.Context, redis_key string, key string) error {
	rdb := GetRedisClient()
	return rdb.HDel(ctx, redis_key, key).Err()
}

func HGet(ctx context.Context, redis_key string, key string) (string, error) {
	rdb := GetRedisClient()
	return rdb.HGet(ctx, redis_key, key).Result()
}

func HGetAll(ctx context.Context, redis_key string) (map[string]string, error) {
	rdb := GetRedisClient()
	return rdb.HGetAll(ctx, redis_key).Result()
}
