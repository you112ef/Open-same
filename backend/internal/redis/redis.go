package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/open-same/backend/internal/config"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

// Init initializes the Redis connection
func Init(cfg config.RedisConfig) (*redis.Client, error) {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: 20,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	log.Println("Redis connection established successfully")
	return Client, nil
}

// GetClient returns the Redis client instance
func GetClient() *redis.Client {
	return Client
}

// Close closes the Redis connection
func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// Set sets a key-value pair with expiration
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return Client.Set(ctx, key, value, expiration).Err()
}

// Get gets a value by key
func Get(ctx context.Context, key string) (string, error) {
	return Client.Get(ctx, key).Result()
}

// GetBytes gets a value as bytes by key
func GetBytes(ctx context.Context, key string) ([]byte, error) {
	return Client.Get(ctx, key).Bytes()
}

// Del deletes a key
func Del(ctx context.Context, keys ...string) error {
	return Client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func Exists(ctx context.Context, key string) (bool, error) {
	result, err := Client.Exists(ctx, key).Result()
	return result > 0, err
}

// Expire sets expiration for a key
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return Client.Expire(ctx, key, expiration).Err()
}

// TTL gets the time to live for a key
func TTL(ctx context.Context, key string) (time.Duration, error) {
	return Client.TTL(ctx, key).Result()
}

// Incr increments a key
func Incr(ctx context.Context, key string) (int64, error) {
	return Client.Incr(ctx, key).Result()
}

// IncrBy increments a key by a specific amount
func IncrBy(ctx context.Context, key, amount int64) (int64, error) {
	return Client.IncrBy(ctx, key, amount).Result()
}

// HSet sets a hash field
func HSet(ctx context.Context, key string, values ...interface{}) error {
	return Client.HSet(ctx, key, values...).Err()
}

// HGet gets a hash field
func HGet(ctx context.Context, key, field string) (string, error) {
	return Client.HGet(ctx, key, field).Result()
}

// HGetAll gets all hash fields
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return Client.HGetAll(ctx, key).Result()
}

// HDel deletes hash fields
func HDel(ctx context.Context, key string, fields ...string) error {
	return Client.HDel(ctx, key, fields...).Err()
}

// SAdd adds members to a set
func SAdd(ctx context.Context, key string, members ...interface{}) error {
	return Client.SAdd(ctx, key, members...).Err()
}

// SMembers gets all members of a set
func SMembers(ctx context.Context, key string) ([]string, error) {
	return Client.SMembers(ctx, key).Result()
}

// SRem removes members from a set
func SRem(ctx context.Context, key string, members ...interface{}) error {
	return Client.SRem(ctx, key, members...).Err()
}

// SIsMember checks if a member is in a set
func SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return Client.SIsMember(ctx, key, member).Result()
}

// ZAdd adds members to a sorted set
func ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return Client.ZAdd(ctx, key, members...).Err()
}

// ZRange gets members from a sorted set by rank
func ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return Client.ZRange(ctx, key, start, stop).Result()
}

// ZRem removes members from a sorted set
func ZRem(ctx context.Context, key string, members ...interface{}) error {
	return Client.ZRem(ctx, key, members...).Err()
}

// Pipeline returns a new pipeline
func Pipeline() redis.Pipeliner {
	return Client.Pipeline()
}

// TxPipeline returns a new transaction pipeline
func TxPipeline() redis.Pipeliner {
	return Client.TxPipeline()
}

// Watch watches keys for changes
func Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	return Client.Watch(ctx, fn, keys...)
}

// Subscribe subscribes to a channel
func Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return Client.Subscribe(ctx, channels...)
}

// Publish publishes a message to a channel
func Publish(ctx context.Context, channel string, message interface{}) error {
	return Client.Publish(ctx, channel, message).Err()
}

// FlushDB flushes the current database
func FlushDB(ctx context.Context) error {
	return Client.FlushDB(ctx).Err()
}

// FlushAll flushes all databases
func FlushAll(ctx context.Context) error {
	return Client.FlushAll(ctx).Err()
}

// Info gets Redis server information
func Info(ctx context.Context, section ...string) (string, error) {
	return Client.Info(ctx, section...).Result()
}

// ClientList gets client connections
func ClientList(ctx context.Context) (string, error) {
	return Client.ClientList(ctx).Result()
}

// MemoryUsage gets memory usage for a key
func MemoryUsage(ctx context.Context, key string, samples ...int) (int64, error) {
	return Client.MemoryUsage(ctx, key, samples...).Result()
}