package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/opencode/erp-school-backend/internal/config"
	"go.uber.org/zap"
)

type MinioClient struct {
	client *minio.Client
	bucket string
	logger *zap.Logger
}

type UploadResult struct {
	URL      string `json:"url"`
	Key      string `json:"key"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

func NewMinioClient(cfg config.MinIOConfig, logger *zap.Logger) (*MinioClient, error) {
	if !cfg.IsConfigured() {
		logger.Warn("minio not configured, storage features disabled")
		return &MinioClient{logger: logger}, nil
	}

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	exists, err := client.BucketExists(context.Background(), cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(context.Background(), cfg.Bucket, minio.MakeBucketOptions{
			Region: cfg.Region,
		}); err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	logger.Info("connected to minio",
		zap.String("endpoint", cfg.Endpoint),
		zap.String("bucket", cfg.Bucket),
	)

	return &MinioClient{
		client: client,
		bucket: cfg.Bucket,
		logger: logger,
	}, nil
}

func (m *MinioClient) IsAvailable() bool {
	return m.client != nil
}

func (m *MinioClient) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (*UploadResult, error) {
	if !m.IsAvailable() {
		return nil, fmt.Errorf("minio client not available")
	}

	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	info, err := m.client.PutObject(ctx, m.bucket, key, reader, size, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	url := fmt.Sprintf("/api/v1/files/%s", key)

	return &UploadResult{
		URL:      url,
		Key:      key,
		Size:     info.Size,
		MimeType: contentType,
	}, nil
}

func (m *MinioClient) UploadFile(ctx context.Context, folder string, filename string, reader io.Reader, size int64, contentType string) (*UploadResult, error) {
	ext := filepath.Ext(filename)
	key := fmt.Sprintf("%s/%d%s", folder, time.Now().UnixNano(), ext)
	return m.Upload(ctx, key, reader, size, contentType)
}

func (m *MinioClient) Download(ctx context.Context, key string) (io.ReadCloser, string, error) {
	if !m.IsAvailable() {
		return nil, "", fmt.Errorf("minio client not available")
	}

	obj, err := m.client.GetObject(ctx, m.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", fmt.Errorf("failed to get object: %w", err)
	}

	stat, err := obj.Stat()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get object stat: %w", err)
	}

	return obj, stat.ContentType, nil
}

func (m *MinioClient) Delete(ctx context.Context, key string) error {
	if !m.IsAvailable() {
		return fmt.Errorf("minio client not available")
	}

	err := m.client.RemoveObject(ctx, m.bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

func (m *MinioClient) DeleteMultiple(ctx context.Context, keys []string) error {
	if !m.IsAvailable() {
		return fmt.Errorf("minio client not available")
	}

	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		for _, key := range keys {
			objectsCh <- minio.ObjectInfo{Key: key}
		}
	}()

	for err := range m.client.RemoveObjects(ctx, m.bucket, objectsCh, minio.RemoveObjectsOptions{}) {
		if err.Err != nil {
			return fmt.Errorf("failed to delete object %s: %w", err.ObjectName, err.Err)
		}
	}
	return nil
}

func (m *MinioClient) PresignedGetURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	if !m.IsAvailable() {
		return "", fmt.Errorf("minio client not available")
	}

	url, err := m.client.PresignedGetObject(ctx, m.bucket, key, expires, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url: %w", err)
	}
	return url.String(), nil
}

func (m *MinioClient) PresignedPutURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	if !m.IsAvailable() {
		return "", fmt.Errorf("minio client not available")
	}

	url, err := m.client.PresignedPutObject(ctx, m.bucket, key, expires)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned put url: %w", err)
	}
	return url.String(), nil
}

func (m *MinioClient) GetObjectURL(key string) string {
	return fmt.Sprintf("/api/v1/files/%s", key)
}
