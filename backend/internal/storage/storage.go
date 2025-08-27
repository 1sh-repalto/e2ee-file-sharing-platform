package storage

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, bucket, objectName string, reader io.Reader, objectSize int64, contentType string) error
	Download(ctx context.Context, bucket, objectName string) (io.ReadCloser, error)
	Delete(ctx context.Context, bucket, objectName string) error
}