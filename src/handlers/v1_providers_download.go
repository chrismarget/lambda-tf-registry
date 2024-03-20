package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/chrismarget/lambda-tf-registry/src/env"
	"github.com/chrismarget/lambda-tf-registry/src/responders"
	"github.com/chrismarget/lambda-tf-registry/src/utils"
	"net/http"
	"strings"
)

// /v1/providers/hashicorp/random/2.0.0/download/linux/amd64

var _ Handler = new(V1ProviderDownloadHandler)

type V1ProviderDownloadHandler struct {
	urlParts []string
}

func (o V1ProviderDownloadHandler) Handle(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	tableName := env.Get(env.ProviderTableName)

	var r responders.Responder
	switch {
	case tableName == "":
		r = responders.Error{Code: http.StatusInternalServerError, Err: fmt.Errorf("env var %q not set", env.ProviderTableName), Req: &req}
	case len(o.urlParts[2]) == 0:
		r = responders.Error{Code: http.StatusUnprocessableEntity, Err: errors.New("url part 2 (namespace) must not be empty"), Req: &req}
	case len(o.urlParts[3]) == 0:
		r = responders.Error{Code: http.StatusUnprocessableEntity, Err: errors.New("url part 3 (type) must not be empty"), Req: &req}
	case len(o.urlParts[4]) == 0:
		r = responders.Error{Code: http.StatusUnprocessableEntity, Err: errors.New("url part 4 (version) must not be empty"), Req: &req}
	case len(o.urlParts[6]) == 0:
		r = responders.Error{Code: http.StatusUnprocessableEntity, Err: errors.New("url part 6 (os) must not be empty"), Req: &req}
	case len(o.urlParts[7]) == 0:
		r = responders.Error{Code: http.StatusUnprocessableEntity, Err: errors.New("url part 7 (arch) must not be empty"), Req: &req}
	}
	if r != nil {
		return r.Respond()
	}

	client, err := utils.DdbClient()
	if err != nil {
		r = responders.Error{Code: http.StatusInternalServerError, Err: err, Req: &req}
		return r.Respond()
	}

	namespaceType := strings.Join(o.urlParts[2:4], "/")
	versionOsArch := strings.Join(append(o.urlParts[4:5], o.urlParts[6:8]...), "/")

	getItemOutput, err := client.GetItem(&dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"NamespaceType": {S: &namespaceType},
			"VersionOsArch": {S: &versionOsArch},
		},
	})
	if err != nil {
		r = responders.Error{Code: http.StatusInternalServerError, Err: err, Req: &req}
		return r.Respond()
	}
	if getItemOutput.Item == nil {
		r = responders.Error{Code: http.StatusNotFound, Err: errors.New(req.RawPath), Req: &req}
		return r.Respond()
	}

	download := responders.V1Download{
		ItemMap:       getItemOutput.Item,
		NamespaceType: namespaceType,
		VersionOsArch: versionOsArch,
	}

	return download.Respond()
}

func newV1ProviderDownloadHandler(urlParts []string) Handler {
	return V1ProviderDownloadHandler{urlParts: urlParts}
}
