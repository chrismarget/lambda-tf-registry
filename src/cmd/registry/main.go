package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/chrismarget/lambda-tf-registry/src/env"
	"github.com/chrismarget/lambda-tf-registry/src/errors"
	"github.com/chrismarget/lambda-tf-registry/src/url"
)

func HandleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	path, err := url.NewPathFromString(request.RawPath)
	if err != nil {
		var ie ierrors.IErr
		if errors.As(err, &ie) {
			return ie.LambdaResponse()
		}
		return events.LambdaFunctionURLResponse{}, err
	}

	sess, err := session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable})
	if err != nil {
		var ie ierrors.IErr
		if errors.As(err, &ie) {
			return ie.LambdaResponse()
		}
		return events.LambdaFunctionURLResponse{}, err
	}

	return path.HandleRequest(dynamodb.New(sess), env.ParseEnv())
}

func main() {
	lambda.Start(HandleRequest)
}
