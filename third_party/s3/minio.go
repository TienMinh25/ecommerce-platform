package s3

import (
	"context"
	"fmt"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"log"
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

func (s *storage) Upload(ctx context.Context, payload pkg.UploadInput, bucket string) (string, error) {
	ctx, span := s.tracer.StartFromContext(ctx, "Upload")
	defer span.End()
}

func (s *storage) Delete(ctx context.Context, name string) error {
	//TODO implement me
	panic("implement me")
}

func (s *storage) GenerateUploadPresignedURL(ctx context.Context, objectName, bucket string) (string, error) {
	//TODO implement me
	panic("implement me")
}
