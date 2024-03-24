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
	"github.com/chrismarget/lambda-tf-registry/src/common"
	"github.com/chrismarget/lambda-tf-registry/src/common/awsclients"
	httpError "github.com/chrismarget/lambda-tf-registry/src/error"
	"github.com/chrismarget/lambda-tf-registry/src/v1_handlers/env"
	"github.com/chrismarget/lambda-tf-registry/src/v1_responses"
)

const providerDownloadPath = "/v1/providers/[^/]+/[^/]+/[0-9.]+/download/[^/]+/[^/]+"

var _ Handler = new(ProviderDownloadHandler)

type ProviderDownloadHandler struct{}

func (o ProviderDownloadHandler) AddRoutes(router *lmdrouter.Router) {
	router.Route(http.MethodGet, providerDownloadPath, o.Handle)
}

func (o ProviderDownloadHandler) Handle(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("%s: %q\n", req.HTTPMethod, req.Path)
	urlParts := strings.Split(strings.TrimLeft(req.Path, common.PathSep), common.PathSep)
	if len(urlParts) != 8 {
		hErr := httpError.FromPrivateError(
			fmt.Errorf("expected URL to have 8 parts, got %q", req.Path),
			"unexpected URL path length",
		)
		return hErr.MarshalResponse()
	}

	namespaceType := strings.Join(urlParts[2:4], common.PathSep)
	versionOsArch := strings.Join(append(urlParts[4:5], urlParts[6:8]...), common.PathSep)

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

	getItemOutput, err := client.GetItem(&dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"NamespaceType": {S: &namespaceType},
			"VersionOsArch": {S: &versionOsArch},
		},
	})
	if err != nil {
		hErr := httpError.FromPrivateError(err, "cannot get item from database")
		return hErr.MarshalResponse()
	}
	if getItemOutput.Item == nil {
		return lmdrouter.MarshalResponse(http.StatusNotFound, nil, nil)
	}

	data := v1responses.Download{
		ItemMap:       getItemOutput.Item,
		NamespaceType: namespaceType,
		VersionOsArch: versionOsArch,
	}

	return lmdrouter.MarshalResponse(http.StatusOK, nil, &data)
}

func NewProviderDownloadHandler() Handler {
	return ProviderDownloadHandler{}
}
