package pkg

import (
	"context"
	"time"
)

type ICache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	GetAndDelete(ctx context.Context, key string) (string, error)

	SetHash(ctx context.Context, key string, data map[string]interface{}, ttl time.Duration) error
	GetHash(ctx context.Context, key string) (map[string]string, error)
	DeleteHash(ctx context.Context, key string) error
	GetAndDeleteHash(ctx context.Context, key string) (map[string]string, error)
}
