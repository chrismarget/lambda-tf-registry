package awsclients

import "github.com/aws/aws-sdk-go/service/dynamodb"

func DdbClient() (*dynamodb.DynamoDB, error) {
	session, err := getSession()
	if err != nil {
		return nil, err
	}

	return dynamodb.New(session), nil
}
