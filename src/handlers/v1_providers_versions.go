package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/chrismarget/lambda-tf-registry/src/env"
	"github.com/chrismarget/lambda-tf-registry/src/responders"
	"net/http"
	"strings"
)

// /v1/providers/hashicorp/random/versions

var _ Handler = new(V1ProviderVersionsHandler)

type V1ProviderVersionsHandler struct {
	urlParts []string
}

func (o V1ProviderVersionsHandler) Handle(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	tableName := env.Get(env.ProviderTableName)

	var r responders.Responder
	switch {
	case tableName == "":
		r = responders.Error{Code: http.StatusInternalServerError, Err: fmt.Errorf("env var %q not set", env.ProviderTableName), Req: &req}
	case o.urlParts[4] != "versions":
		r = responders.Error{Code: http.StatusNotFound, Req: &req}
	case len(o.urlParts[2]) == 0:
		r = responders.Error{Code: http.StatusUnprocessableEntity, Err: errors.New("url part 2 (namespace) must not be empty"), Req: &req}
	case len(o.urlParts[3]) == 0:
		r = responders.Error{Code: http.StatusUnprocessableEntity, Err: errors.New("url part 3 (type) must not be empty"), Req: &req}
	}
	if r != nil {
		return r.Respond()
	}

	client, err := ddbClient()
	if err != nil {
		r = responders.Error{Code: http.StatusInternalServerError, Err: err, Req: &req}
		return r.Respond()
	}

	// aws dynamodb query --table-name registry-providers --projection-expression SHA --key-condition-expression "NamespaceType = :v1" --expression-attribute-values '{":v1":{"S":"hashicorp/tls"}}'
	namespaceType := strings.Join(o.urlParts[2:4], "/")

	keyCondition := expression.Key("NamespaceType").Equal(expression.Value(namespaceType))
	projection := expression.NamesList(
		expression.Name("VersionOsArch"),
		//expression.Name("Keys"),
		expression.Name("Protocols"),
		//expression.Name("SHA"),
		//expression.Name("SHA_URL"),
		//expression.Name("Sig_URL"),
		//expression.Name("URL"),
	)
	expr, err := expression.NewBuilder().
		WithKeyCondition(keyCondition).
		WithProjection(projection).
		Build()
	if err != nil {
		r = responders.Error{Code: http.StatusInternalServerError, Err: err, Req: &req}
		return r.Respond()
	}

	queryOutput, err := client.Query(&dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 &tableName,
	})
	if err != nil {
		r = responders.Error{Code: http.StatusInternalServerError, Err: err, Req: &req}
		return r.Respond()
	}
	if len(queryOutput.Items) == 0 {
		r = responders.Error{Code: http.StatusNotFound, Req: &req}
		return r.Respond()
	}

	versions := responders.V1Versions{
		Items:         queryOutput.Items,
		NamespaceType: namespaceType,
	}

	return versions.Respond()
}

func newV1ProviderVersionsHandler(urlParts []string) Handler {
	return V1ProviderVersionsHandler{urlParts: urlParts}
}
