package handlers

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

var _ Handler = new(V1ProviderHandler)

type V1ProviderHandler struct {
	urlParts []string
}

func (o V1ProviderHandler) Handle(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	var h Handler
	switch len(o.urlParts) {
	case 5:
		h = newV1ProviderVersionsHandler(o.urlParts)
	case 8:
		h = newV1ProviderDownloadHandler(o.urlParts)
	default:
		h = newErrorHandler(http.StatusNotFound, nil)
	}

	return h.Handle(ctx, req)
}

func newV1ProviderHandler(urlParts []string) Handler {
	return V1ProviderHandler{urlParts: urlParts}
}
