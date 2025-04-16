package pkg

import (
	"context"
	"io"
)

type Storage interface {
	// Upload uploads a file to the specified bucket (or default bucket if empty)
	// and returns the file's URL.
	// (e.g., file content, metadata) and returns the file's unique identifier or URL.
	Upload(ctx context.Context, payload UploadInput, bucket string) (string, error)

	// Delete deletes a file from the specified bucket (or default bucket if empty)
	Delete(ctx context.Context, name string) error

	// GenerateUploadPresignedURL generates a pre-signed URL for uploading to the
	// specified bucket (or default bucket if empty)
	GenerateUploadPresignedURL(ctx context.Context, objectName, bucket string) (string, error)
}

type UploadInput struct {
	File        io.Reader
	Name        string
	Size        int64
	ContentType string
}
