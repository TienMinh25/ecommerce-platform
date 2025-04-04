package api_gateway_service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"time"
)

type oauthCacheService struct {
	redisCache pkg.ICache
	tracer     pkg.Tracer
}

func NewOauthCacheService(redisCache pkg.ICache, tracer pkg.Tracer) IOauthCacheService {
	return &oauthCacheService{
		redisCache: redisCache,
		tracer:     tracer,
	}
}

func (o *oauthCacheService) SaveOauthState(ctx context.Context, state string) error {
	ctx, span := o.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "SaveOauthState"))
	defer span.End()

	return o.redisCache.Set(ctx, "oauth:state:"+state, state, time.Minute*7)
}

func (o *oauthCacheService) GetAndDeleteOauthState(ctx context.Context, state string) (string, error) {
	ctx, span := o.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetAndDeleteOauthState"))
	defer span.End()

	stateFromCache, err := o.redisCache.GetAndDelete(ctx, "oauth:state:"+state)

	if err != nil {
		return "", err
	}

	return stateFromCache, nil
}
