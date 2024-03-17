package url

import (
	"fmt"
	"net/http"

	"github.com/chrismarget/lambda-tf-registry/src/env"
	ierrors "github.com/chrismarget/lambda-tf-registry/src/errors"
	"github.com/chrismarget/lambda-tf-registry/src/response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	providerDownloadPathAction    = "download"
	providerDownloadPathPartCount = 8
)

var _ Path = new(ProviderDownloadPath)

type ProviderDownloadPath struct {
	apiVersion      string
	namespace       string
	providerType    string
	providerVersion string
	os              string
	arch            string
}

func (o *ProviderDownloadPath) loadParts(parts []string) (Path, error) {
	if len(parts) != providerDownloadPathPartCount {
		return nil, newPathError(
			http.StatusUnprocessableEntity,
			"provider %q URL must have %d parts, got %s",
			providerDownloadPathAction,
			providerDownloadPathPartCount,
			parts)
	}

	if parts[5] != providerDownloadPathAction {
		return nil, newPathError(
			http.StatusUnprocessableEntity,
			"part [5] of provider URL must be %q, got %q",
			providerDownloadPathAction,
			parts[5],
		)
	}

	o.apiVersion = parts[0]
	o.namespace = parts[2]
	o.providerType = parts[3]
	o.providerVersion = parts[4]
	o.os = parts[6]
	o.arch = parts[7]

	return o, nil
}

func (o *ProviderDownloadPath) ApiVersion() string {
	return o.apiVersion
}

func (o *ProviderDownloadPath) HandleRequest(ddb *dynamodb.DynamoDB, env env.Env) (events.LambdaFunctionURLResponse, error) {
	gio, err := ddb.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(env.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"NamespaceType": {
				S: o.NamespaceType(),
			},
			"VersionOsArch": {
				S: o.VersionOsArch(),
			},
		},
	})
	if err != nil {
		return events.LambdaFunctionURLResponse{}, err
	}
	if gio.Item == nil {
		return ierrors.IErr{
			Err:  fmt.Errorf("not found: %s/%s", *o.NamespaceType(), *o.VersionOsArch()),
			Code: http.StatusNotFound,
		}.LambdaResponse()
	}

	download := response.Download{
		NamespaceType: *o.NamespaceType(),
		VersionOsArch: *o.VersionOsArch(),
	}

	return download.Respond(gio.Item)
}

func (o *ProviderDownloadPath) NamespaceType() *string {
	s := fmt.Sprintf("%s/%s", o.namespace, o.providerType)
	return &s
}

func (o *ProviderDownloadPath) Type() PathType {
	return PathTypeProvider
}

func (o *ProviderDownloadPath) VersionOsArch() *string {
	s := fmt.Sprintf("%s/%s/%s", o.providerVersion, o.os, o.arch)
	return &s
}
