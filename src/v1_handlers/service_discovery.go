package v1handlers

import (
	"context"
	"net/http"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
)

const serviceDiscoveryPath = "/.well-known/terraform.json"

var _ Handler = new(ServiceDiscoveryHandler)

type ServiceDiscoveryHandler struct{ services map[string]string }

func (o ServiceDiscoveryHandler) AddRoutes(router *lmdrouter.Router) {
	router.Route("GET", serviceDiscoveryPath, o.Handle)
}

func (o ServiceDiscoveryHandler) Handle(_ context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return lmdrouter.MarshalResponse(http.StatusOK, nil, o.services)
}

func NewServiceDiscoveryHandler(services map[string]string) Handler {
	return ServiceDiscoveryHandler{
		services: services,
	}
}
