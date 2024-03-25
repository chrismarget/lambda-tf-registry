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

const providerDownloadPath = "/v1/providers/[^/]+/[^/]+/[0-9.]+/download/[^/]+/[^/]+"

var _ Handler = new(ProviderDownloadHandler)

type ProviderDownloadHandler struct{}

func (o ProviderDownloadHandler) AddRoutes(router *lmdrouter.Router) {
	router.Route(http.MethodGet, providerDownloadPath, o.Handle)
}

func (o ProviderDownloadHandler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("%s: %q\n", req.HTTPMethod, req.Path)

	tableName := common.Get(common.ProviderTableName)
	if tableName == "" {
		hErr := common.NewPublicError("cannot determine database table name")
		return hErr.MarshalResponse()
	}

	client, err := awsclients.DdbClient(ctx)
	if err != nil {
		hErr := common.FromPrivateError(err, "cannot create database client")
		return hErr.MarshalResponse()
	}

	model, err := NewDownloadModelFromUrlPath(req.Path)
	if err != nil {
		hErr := common.FromPrivateError(err, "failed creating db query from URL path")
		return hErr.MarshalResponse()
	}

	queryKey, err := model.GetKey()
	if err != nil {
		hErr := common.NewPublicError("failed to marshal query key")
		return hErr.MarshalResponse()
	}

	response, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       queryKey,
		TableName: &tableName,
	})
	if err != nil {
		hErr := common.FromPrivateError(err, "database error")
		return hErr.MarshalResponse()
	}

	err = attributevalue.UnmarshalMap(response.Item, &model)
	if err != nil {
		hErr := common.FromPrivateError(err, "database unmarshal error")
		return hErr.MarshalResponse()
	}

	return lmdrouter.MarshalResponse(http.StatusOK, nil, &model)
}

func NewProviderDownloadHandler() Handler {
	return ProviderDownloadHandler{}
}
