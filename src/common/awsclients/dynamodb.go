package awsclients

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func DdbClient(ctx context.Context) (*dynamodb.Client, error) {
	cfg, err := loadConfig(ctx)
	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(*cfg), nil
}
