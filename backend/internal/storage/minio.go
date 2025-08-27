package storage

import (
	"context"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client *minio.Client
}

func NewMinioStorage(endpoint, accessKey, secretKey string, useSSL bool) (*MinioStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &MinioStorage{client: client}, nil
}

func (s *MinioStorage) Upload(ctx context.Context, bucket, objectName string, reader io.Reader, objectSize int64, contentType string) error {
	exists, err := s.client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}

	if !exists {
		if err := s.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return err
		}
		log.Printf("Created bucket %s\n", bucket)
	}

	_, err = s.client.PutObject(ctx, bucket, objectName, reader, objectSize, minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (s *MinioStorage) Download(ctx context.Context, bucket, objectName string) (io.ReadCloser, error) {
	return s.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
}

func (s *MinioStorage) Delete(ctx context.Context, bucket, objectName string) error {
	return s.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
}