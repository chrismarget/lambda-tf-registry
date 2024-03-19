package handlers

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func ddbClient() (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable})
	if err != nil {
		return nil, err
	}

	return dynamodb.New(sess), nil
}
