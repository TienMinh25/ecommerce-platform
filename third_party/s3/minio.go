package s3

import (
	"context"
	"fmt"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strings"
	"time"
)

type storage struct {
	client        *minio.Client
	defaultBucket string
	region        string
	tracer        pkg.Tracer
	endpointURL   string
}

func NewStorage(env *env.EnvManager, tracer pkg.Tracer) (pkg.Storage, error) {
	useSSL := false

	minioClient, err := minio.New(env.MinioConfig.MinioEndpointURL, &minio.Options{
		Creds:  credentials.NewStaticV4(env.MinioConfig.MinioAccessKey, env.MinioConfig.MinioSecretKey, ""),
		Secure: useSSL,
	})

	fmt.Println("✅ Connected to S3 successfully!")

	if err != nil {
		return nil, errors.Wrap(err, "minio.New")
	}

	ctx := context.Background()

	// Set default bucket
	defaultBucket := env.MinioConfig.MinioBucketAvatars

	// Check if default bucket exists and create if needed
	isBucketExists, err := minioClient.BucketExists(ctx, defaultBucket)

	if err == nil && !isBucketExists {
		err = minioClient.MakeBucket(ctx, env.MinioConfig.MinioBucketAvatars, minio.MakeBucketOptions{
			Region: env.MinioConfig.MinioRegion,
		})

		if err != nil {
			return nil, errors.Wrap(err, "minioClient.MakeBucket")
		}

		log.Println("policy: %v", env.MinioConfig.MinioBucketAvatarsPolicy)
		err = minioClient.SetBucketPolicy(ctx, env.MinioConfig.MinioBucketAvatars, env.MinioConfig.MinioBucketAvatarsPolicy)

		if err != nil {
			return nil, errors.Wrap(err, "minioClient.SetBucketPolicy")
		}
	}

	return &storage{
		client:        minioClient,
		defaultBucket: defaultBucket,
		region:        env.MinioConfig.MinioRegion,
		tracer:        tracer,
		endpointURL:   env.MinioConfig.MinioEndpointURL,
	}, nil
}

// getBucketName returns the specified bucket or the default if none is specified
func (s *storage) getBucketName(bucket string) string {
	if bucket == "" {
		return s.defaultBucket
	}

	return bucket
}

// getURLPrefix generates URL prefix for the given bucket
func (s *storage) getURLPrefix(bucket string) string {
	endpointURL := s.client.EndpointURL()

	return fmt.Sprintf("%s://%s.s3.%s.%s/", endpointURL.Scheme, bucket, s.region, endpointURL.Host)
}

// Upload implements pkg.Storage interface with bucket support
func (s *storage) Upload(ctx context.Context, payload pkg.UploadInput, bucket string) (string, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.InfraLayer, "minio.Upload"))
	defer span.End()

	bucketName := s.getBucketName(bucket)

	// ensure bucket exists
	exists, err := s.client.BucketExists(ctx, bucketName)

	if err != nil {
		return "", utils.TechnicalError{
			Message: fmt.Sprintf("Failed to check if bucket exists: %v", err),
			Code:    http.StatusInternalServerError,
		}
	}

	// create bucket if it doesn't exist
	if !exists {
		err = s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region: s.region,
		})

		if err != nil {
			return "", utils.TechnicalError{
				Message: fmt.Sprintf("Failed to create bucket: %v", err),
				Code:    http.StatusInternalServerError,
			}
		}
	}

	// Upload the object
	info, err := s.client.PutObject(ctx, bucketName, payload.Name, payload.File, payload.Size, minio.PutObjectOptions{
		ContentType: payload.ContentType,
	})

	if err != nil {
		return "", utils.TechnicalError{
			Message: fmt.Sprintf("Failed to upload: %v", err),
			Code:    http.StatusInternalServerError,
		}
	}

	return s.getURLPrefix(bucketName) + info.Key, nil
}

// Delete implements pkg.Storage interface with bucket support
func (s *storage) Delete(ctx context.Context, name, bucket string) error {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.InfraLayer, "minio.Delete"))
	defer span.End()

	bucketName := s.getBucketName(bucket)
	objectKey := name

	if strings.HasPrefix(objectKey, "/") {
		objectKey = objectKey[1:]
	}

	if err := s.client.RemoveObject(ctx, bucketName, objectKey, minio.RemoveObjectOptions{}); err != nil {
		return utils.TechnicalError{
			Message: fmt.Sprintf("Failed to delete object: %v", err),
			Code:    http.StatusInternalServerError,
		}
	}

	return nil
}

// GenerateUploadPresignedURL implements pkg.Storage interface with bucket support
func (s *storage) GenerateUploadPresignedURL(ctx context.Context, objectName, bucket string) (string, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.InfraLayer, "minio.GenerateUploadPresignedURL"))
	defer span.End()

	bucketName := s.getBucketName(bucket)

	exists, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		return "", utils.TechnicalError{
			Message: fmt.Sprintf("Failed to check if bucket exists: %v", err),
			Code:    http.StatusInternalServerError,
		}
	}

	if !exists {
		err = s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region: s.region,
		})
		if err != nil {
			return "", utils.TechnicalError{
				Message: fmt.Sprintf("Failed to create bucket: %v", err),
				Code:    http.StatusInternalServerError,
			}
		}
	}

	// generate presigned url
	presignedURL, err := s.client.PresignedPutObject(ctx, bucketName, objectName, time.Hour)

	if err != nil {
		return "", utils.TechnicalError{
			Message: fmt.Sprintf("Failed to generate presigned URL: %v", err),
			Code:    http.StatusInternalServerError,
		}
	}

	return presignedURL.String(), nil
}
