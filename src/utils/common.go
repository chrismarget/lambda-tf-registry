package utils

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)


var sess *session.Session

func getSession() (*session.Session, error) {
	var err error
	if sess == nil {
		sess, err = session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable})
	}

	return sess, err
}

func DdbClient() (*dynamodb.DynamoDB, error) {
	var err error
	sess, err = getSession()
	if err != nil {
		return nil, err
	}

	return dynamodb.New(sess), nil
}

func SmClient() (*secretsmanager.SecretsManager, error) {
	var err error
	sess, err = getSession()
	if err != nil {
		return nil, err
	}

	return secretsmanager.New(sess), nil
}


