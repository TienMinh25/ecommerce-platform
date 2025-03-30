package api_gateway_service

import (
	"context"
	"errors"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/redis/go-redis/v9"
	"net/http"
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

func (o *otpCache) CacheOTP(ctx context.Context, otp, email string, ttl time.Duration) error {
	ctx, span := o.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CacheOTP"))
	defer span.End()

	return o.redisCache.Set(ctx, otp, email, ttl)
}

func (o *otpCache) GetValueString(ctx context.Context, key string) (string, error) {
	ctx, span := o.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetValueString"))
	defer span.End()

	res, err := o.redisCache.Get(ctx, key)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", utils.BusinessError{
				Code:      http.StatusBadRequest,
				Message:   "OTP has expired, please request a new one",
				ErrorCode: errorcode.OTP_TIMEOUT,
			}
		}

		return "", utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return res, nil
}

func (o *otpCache) DeleteOTP(ctx context.Context, otp string) error {
	ctx, span := o.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "DeleteOTP"))
	defer span.End()

	return o.redisCache.Delete(ctx, otp)
}
