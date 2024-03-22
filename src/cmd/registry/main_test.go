package main

import (
	"context"
	"log"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
)

func TestThing(t *testing.T) {
	ctx := context.Background()

	req := events.APIGatewayProxyRequest{
		Resource:                        "",
		Path:                            "/v1/providers/hashicorp/tls/4.0.1/download/linux/amd64",
		HTTPMethod:                      "GET",
		Headers:                         nil,
		MultiValueHeaders:               nil,
		QueryStringParameters:           nil,
		MultiValueQueryStringParameters: nil,
		PathParameters:                  nil,
		StageVariables:                  nil,
		RequestContext:                  events.APIGatewayProxyRequestContext{},
		Body:                            "",
		IsBase64Encoded:                 false,
	}
	resp, err := router.Handler(ctx, req)
	require.NoError(t, err)
	log.Println(resp.Body)
}
