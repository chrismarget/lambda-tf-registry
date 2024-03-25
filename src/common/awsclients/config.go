package awsclients

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var awsConfig *aws.Config

func loadConfig(ctx context.Context) (*aws.Config, error) {
	if awsConfig != nil {
		return awsConfig, nil
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	awsConfig = &cfg

	return &cfg, nil
}
