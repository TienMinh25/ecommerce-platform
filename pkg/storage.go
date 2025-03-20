package pkg

import (
	"context"
	"io"
)

type Storage interface {
	// Upload Uploads a file or data to the storage system using the
	// provided `UploadInput`
	// (e.g., file content, metadata) and returns the file's unique identifier or URL.
	Upload(ctx context.Context, payload UploadInput) (string, error)

	// Delete Deletes a file or object from the storage system by its name or identifier, returning an error if the operation fails
	Delete(ctx context.Context, name string) error
}

type UploadInput struct {
	File        io.Reader
	Name        string
	Size        int64
	ContentType string
}
