package v1handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/chrismarget/lambda-tf-registry/src/common"
	"github.com/chrismarget/lambda-tf-registry/src/common/awsclients"
	httpError "github.com/chrismarget/lambda-tf-registry/src/error"
	"github.com/chrismarget/lambda-tf-registry/src/v1_handlers/env"
	v1responses "github.com/chrismarget/lambda-tf-registry/src/v1_responses"
)

const providerVersionsPath = "/v1/providers/[^/]+/[^/]+/versions"

var _ Handler = new(ProviderVersionsHandler)

type ProviderVersionsHandler struct{}

func (o ProviderVersionsHandler) AddRoutes(router *lmdrouter.Router) {
	router.Route(http.MethodGet, providerVersionsPath, o.Handle)
}

func (o ProviderVersionsHandler) Handle(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("%s: %q\n", req.HTTPMethod, req.Path)
	urlParts := strings.Split(strings.TrimLeft(req.Path, common.PathSep), common.PathSep)
	if len(urlParts) != 5 {
		hErr := httpError.FromPrivateError(
			fmt.Errorf("expected URL to have 5 parts, got %q", req.Path),
			"unexpected URL path length",
		)
		return hErr.MarshalResponse()
	}

	namespaceType := strings.Join(urlParts[2:4], common.PathSep)

	tableName := env.Get(env.ProviderTableName)
	if tableName == "" {
		hErr := httpError.NewPublicError("cannot determine database table name")
		return hErr.MarshalResponse()
	}

	client, err := awsclients.DdbClient()
	if err != nil {
		hErr := httpError.FromPrivateError(err, "cannot create database client")
		return hErr.MarshalResponse()
	}

	keyCondition := expression.Key("NamespaceType").Equal(expression.Value(namespaceType))
	projection := expression.NamesList(
		expression.Name("VersionOsArch"),
		expression.Name("Protocols"),
	)
	expr, err := expression.NewBuilder().
		WithKeyCondition(keyCondition).
		WithProjection(projection).
		Build()
	if err != nil {
		hErr := httpError.FromPrivateError(err, "failed building database expression")
		return hErr.MarshalResponse()
	}

	queryOutput, err := client.Query(&dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 &tableName,
	})
	if err != nil {
		hErr := httpError.FromPrivateError(err, "failed querying database")
		return hErr.MarshalResponse()
	}
	if len(queryOutput.Items) == 0 {
		return lmdrouter.MarshalResponse(http.StatusNotFound, nil, nil)
	}

	data := v1responses.Versions{
		ItemMaps:      queryOutput.Items,
		NamespaceType: namespaceType,
	}

	return lmdrouter.MarshalResponse(http.StatusOK, nil, &data)
}

func NewProviderVersionsHandler() Handler {
	return ProviderVersionsHandler{}
}
