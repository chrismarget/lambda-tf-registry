package v1handlers

import (
	"context"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
)

type Handler interface {
	AddRoutes(*lmdrouter.Router)
	Handle(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}
