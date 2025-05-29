package api_gateway_service

import (
	"context"
	"fmt"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/google/uuid"
	"path/filepath"
	"time"
)

type s3Service struct {
	tracer pkg.Tracer
	minio  pkg.Storage
}

func NewS3Service(tracer pkg.Tracer, minio pkg.Storage) IS3Service {
	return &s3Service{
		tracer: tracer,
		minio:  minio,
	}
}

func (s *s3Service) GetPresignedURLUpload(ctx context.Context, data *api_gateway_dto.GetPresignedURLRequest, userID int) (string, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetPresignedURLUpload"))
	defer span.End()

	// handle to get extension name
	fileExt := filepath.Ext(data.FileName)

	fileUUID := uuid.New().String()
	timestamp := time.Now().UnixNano()

	objectName := fmt.Sprintf("users/%v/%d_%s%s",
		userID,
		timestamp,
		fileUUID,
		fileExt,
	)

	presignedURL, err := s.minio.GenerateUploadPresignedURL(ctx, objectName, string(data.BucketName))

	if err != nil {
		span.RecordError(err)
		return "", err
	}

	return presignedURL, nil
}
