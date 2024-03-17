package url

import (
	"fmt"
	"net/http"

	"github.com/chrismarget/lambda-tf-registry/src/env"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	providerVersionsPathAction    = "versions"
	providerVersionsPathPartCount = 5
)

var _ Path = new(ProviderVersionsPath)

type ProviderVersionsPath struct {
	apiVersion string
	namespace  string
	// original     string
	providerType string
}

//func (o *ProviderVersionsPath) setString(s string) *ProviderVersionsPath {
//	o.original = s
//	return o
//}

func (o *ProviderVersionsPath) loadParts(parts []string) (Path, error) {
	if len(parts) != providerVersionsPathPartCount {
		return nil, newPathError(
			http.StatusUnprocessableEntity,
			"provider %q URL must have %d parts, got %s",
			providerVersionsPathAction,
			providerDownloadPathPartCount,
			parts)
	}

	if parts[4] != providerVersionsPathAction {
		return nil, newPathError(
			http.StatusUnprocessableEntity,
			"part [4] of provider URL must be %q, got %q",
			providerVersionsPathAction,
			parts[4],
		)
	}

	o.apiVersion = parts[0]
	o.namespace = parts[2]
	o.providerType = parts[3]

	return o, nil
}

func (o *ProviderVersionsPath) ApiVersion() string {
	return o.apiVersion
}

func (o *ProviderVersionsPath) HandleRequest(ddb *dynamodb.DynamoDB, env env.Env) (events.LambdaFunctionURLResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (o *ProviderVersionsPath) NamespaceType() *string {
	s := fmt.Sprintf("%s/%s", o.namespace, o.providerType)
	return &s
}

//func (o *ProviderVersionsPath) Original() string {
//	return o.original
//}

func (o *ProviderVersionsPath) Type() PathType {
	return PathTypeProvider
}

func (o *ProviderVersionsPath) VersionOsArch() string {
	// TODO implement me
	panic("implement me")
}
