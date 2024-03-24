package awsclients

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var awsV2Config *aws.Config

func loadConfig(ctx context.Context) (*aws.Config, error) {
	if awsV2Config != nil {
		return awsV2Config, nil
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	awsV2Config = &cfg

	return &cfg, nil
}
