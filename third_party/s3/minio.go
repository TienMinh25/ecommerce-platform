package s3

import (
	"context"
	"fmt"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strings"
)

type storage struct {
	client     *minio.Client
	bucketName string
	urlPrefix  string
	region     string
}

func NewStorage(env *env.EnvManager) (pkg.Storage, error) {
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

	isBucketExists, err := minioClient.BucketExists(ctx, env.MinioConfig.MinioBucketAvatars)

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

	endpointURL := minioClient.EndpointURL()
	urlPrefix := endpointURL.Scheme + "://" + env.MinioConfig.MinioBucketAvatars + ".s3." + env.MinioConfig.MinioRegion + "." + endpointURL.Host + "/"

	return &storage{
		client:     minioClient,
		bucketName: env.MinioConfig.MinioBucketAvatars,
		urlPrefix:  urlPrefix,
		region:     env.MinioConfig.MinioRegion}, nil
}

// Delete implements pkg.Storage.
func (s storage) Upload(ctx context.Context, payload pkg.UploadInput) (string, error) {
	info, err := s.client.PutObject(ctx, s.bucketName, payload.Name, payload.File, payload.Size, minio.PutObjectOptions{
		ContentType: payload.ContentType,
	})

	// if err exists, upload to minio fail
	if err != nil {
		return "", utils.TechnicalError{
			Message: "Failed to upload object to S3",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.urlPrefix + info.Key, nil
}

// Upload implements pkg.Storage.
func (s storage) Delete(ctx context.Context, name string) error {
	// get key object (get path of file on s3)
	name = strings.TrimPrefix(name, s.urlPrefix)

	err := s.client.RemoveObject(ctx, s.bucketName, name, minio.RemoveObjectOptions{})

	if err != nil {
		return utils.TechnicalError{
			Message: "Failed to delete object on S3",
			Code:    http.StatusInternalServerError,
		}
	}

	return nil
}
