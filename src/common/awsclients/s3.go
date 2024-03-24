package awsclients

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func S3Manager(ctx context.Context) (*manager.Uploader, error) {
	s3Client, err := S3Client(ctx)
	if err != nil {
		return nil, err
	}

	return manager.NewUploader(s3Client), nil
}

func S3Client(ctx context.Context) (*s3.Client, error) {
	cfg, err := loadConfig(ctx)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(*cfg), nil
}
