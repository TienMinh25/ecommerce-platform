package api_gateway_service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"time"
)

type otpCache struct {
	redisCache pkg.ICache
	tracer     pkg.Tracer
}

func NewOTPCacheService(redisCache pkg.ICache, tracer pkg.Tracer) IOtpCacheService {
	return &otpCache{
		redisCache: redisCache,
		tracer:     tracer,
	}
}

func (o *otpCache) CacheOTP(ctx context.Context, otp string, ttl time.Duration) error {
	ctx, span := o.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CacheOTP"))
	defer span.End()

	return o.redisCache.Set(ctx, otp, 1, ttl)
}
