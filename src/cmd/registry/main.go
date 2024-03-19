package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chrismarget/lambda-tf-registry/src/handlers"
	"log"
)

func HandleRequest(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	log.Println("enter HandleRequest")
	log.Println("path is: ", req.RawPath)

	h := handlers.NewHandlerFromPath(req.RawPath)
	return h.Handle(ctx, req)
}

func main() {
	log.Println("enter main")
	lambda.Start(HandleRequest)
	log.Println("exit main")
}
