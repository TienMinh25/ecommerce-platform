package infrastructure

import (
	"context"
	"fmt"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type redisCache struct {
	client *redis.Client
	tracer pkg.Tracer
}

func NewRedisCache(env *env.EnvManager, tracer pkg.Tracer) pkg.ICache {
	client := redis.NewClient(&redis.Options{
		Addr:     env.Redis.RedisAddress,
		Password: env.Redis.RedisPassword,
		DB:       env.Redis.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return &redisCache{
		client: client,
		tracer: tracer,
	}
}

// Set key and value string into redis
// Using redis.KeepTTL to pass into ttl param -> keep key and value exists in redis, not be deleted after some time.
func (r *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.InfraLayer, "Set"))
	defer span.End()

	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.InfraLayer, "Get"))
	defer span.End()

	return r.client.Get(ctx, key).Result()
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.InfraLayer, "Delete"))
	defer span.End()

	return r.client.Del(ctx, key).Err()
}

func (r *redisCache) GetAndDelete(ctx context.Context, key string) (string, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.InfraLayer, "GetAndDelete"))
	defer span.End()

	return r.client.GetDel(ctx, key).Result()
}

func (r *redisCache) SetHash(ctx context.Context, key string, data map[string]interface{}, ttl time.Duration) error {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.InfraLayer, "SetHash"))
	defer span.End()

	pipe := r.client.TxPipeline()

	pipe.HSet(ctx, key, data)
	if ttl > 0 {
		pipe.Expire(ctx, key, ttl)
	}

	_, err := pipe.Exec(ctx)

	return err
}

func (r *redisCache) GetHash(ctx context.Context, key string) (map[string]string, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.InfraLayer, "GetHash"))
	defer span.End()

	res, err := r.client.HGetAll(ctx, key).Result()

	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("key '%s' not found", key)
	}

	return res, nil
}

func (r *redisCache) DeleteHash(ctx context.Context, key string) error {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.InfraLayer, "DeleteHash"))
	defer span.End()

	return r.client.Del(ctx, key).Err()
}

func (r *redisCache) GetAndDeleteHash(ctx context.Context, key string) (map[string]string, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.InfraLayer, "GetAndDeleteHash"))
	defer span.End()

	pipe := r.client.TxPipeline()
	getRes := pipe.HGetAll(ctx, key)
	delRes := pipe.Del(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	res, err := getRes.Result()

	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("key '%s' not found", key)
	}

	countEffected, err := delRes.Result()

	if err != nil {
		return nil, err
	}

	if countEffected == 0 {
		return res, fmt.Errorf("no key '%s' in redis to remove", key)
	}

	return res, nil
}
