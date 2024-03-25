package v1handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/chrismarget/lambda-tf-registry/common"
	"github.com/chrismarget/lambda-tf-registry/common/awsclients"
)

const providerVersionsPath = "/v1/providers/[^/]+/[^/]+/versions"

var _ Handler = new(ProviderVersionsHandler)

type ProviderVersionsHandler struct{}

func (o ProviderVersionsHandler) AddRoutes(router *lmdrouter.Router) {
	router.Route(http.MethodGet, providerVersionsPath, o.Handle)
}

func (o ProviderVersionsHandler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("%s: %q\n", req.HTTPMethod, req.Path)

	tableName := common.Get(common.ProviderTableName)
	if tableName == "" {
		hErr := common.NewPublicError("cannot determine database table name")
		return hErr.MarshalResponse()
	}

	model, err := NewVersionsModelFromUrlPath(req.Path)
	if tableName == "" {
		hErr := common.FromPrivateError(err, "cannot determine database table name")
		return hErr.MarshalResponse()
	}

	expr, err := model.KeyExpr()
	if err != nil {
		hErr := common.FromPrivateError(err, "failed to build query")
		return hErr.MarshalResponse()
	}

	client, err := awsclients.DdbClient(ctx)
	if err != nil {
		hErr := common.FromPrivateError(err, "cannot create database client")
		return hErr.MarshalResponse()
	}

	queryPaginator := dynamodb.NewQueryPaginator(client, &dynamodb.QueryInput{
		TableName:                 &tableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})

	for queryPaginator.HasMorePages() {
		response, err := queryPaginator.NextPage(context.TODO())
		if err != nil {
			hErr := common.FromPrivateError(err, "failed querying database")
			return hErr.MarshalResponse()
		}

		var providersPage []ProviderVersionModel
		err = attributevalue.UnmarshalListOfMaps(response.Items, &providersPage)
		if err != nil {
			hErr := common.FromPrivateError(err, "failed unmarshaling database response")
			return hErr.MarshalResponse()
		}

		model.Versions = append(model.Versions, providersPage...)
	}

	return lmdrouter.MarshalResponse(http.StatusOK, nil, &model)
}

func NewProviderVersionsHandler() Handler {
	return ProviderVersionsHandler{}
}
