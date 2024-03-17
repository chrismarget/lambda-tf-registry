package url

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/chrismarget/lambda-tf-registry/src/env"
	ierrors "github.com/chrismarget/lambda-tf-registry/src/errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Path interface {
	ApiVersion() string
	HandleRequest(*dynamodb.DynamoDB, env.Env) (events.LambdaFunctionURLResponse, error)
	NamespaceType() *string
	Type() PathType
}

func NewPathFromString(s string) (Path, error) {
	parts := strings.Split(strings.TrimLeft(s, "/"), "/")
	if len(parts) < 2 {
		return nil, ierrors.IErr{
			Code: http.StatusNotFound,
			Err:  fmt.Errorf("path %q has only %d parts", s, len(parts)),
		}
	}

	if apiVersions.Parse(parts[0]) == nil {
		return nil, ierrors.IErr{
			Code: http.StatusNotFound,
			Err:  fmt.Errorf("path %q indicates unsupported version %q", s, parts[0]),
		}
	}

	pathType := PathTypes.Parse(parts[1])
	if pathType == nil {
		return nil, ierrors.IErr{
			Code: http.StatusNotFound,
			Err:  fmt.Errorf("path %q indicates unsupported type %q", s, parts[1]),
		}
	}

	switch *pathType {
	case PathTypeProvider:
		return parseProviderPath(parts)
	case PathTypeModule:
		return parseModulePath(parts)
	}

	switch {
	case *pathType == PathTypeProvider:
		path := ProviderDownloadPath{apiVersion: parts[0]}
		return path.loadParts(parts)
	}

	return nil, ierrors.IErr{
		Err:  fmt.Errorf("unhandled path type %q", pathType.Value),
		Code: http.StatusNotFound,
	}
}
