package handlers

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/chrismarget/lambda-tf-registry/src/responders"
	"log"
	"net/http"
	"strings"
)

var _ Handler = new(wellKnownHandler)

type wellKnownHandler struct{}

func (o wellKnownHandler) Handle(_ context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	log.Println("enter wellKnownHandler.Handle()")
	urlParts := strings.Split(strings.TrimLeft(req.RawPath, "/"), "/")

	var r responders.Responder
	switch {
	case len(urlParts) != 2:
		r = responders.Error{
			Code: http.StatusNotFound,
			Err:  errors.New("URL path has wrong part count"),
			Req:  &req,
		}
	case urlParts[1] == "terraform.json":
		r = responders.WellKnownTerraform{
			ModulesV1:   "/v1/modules/",
			ProvidersV1: "/v1/providers/",
		}
	default:
		r = responders.Error{
			Code: http.StatusNotFound,
			Req:  &req,
		}
	}

	return r.Respond()
}

func newWellKnownHandler() Handler {
	return wellKnownHandler{}
}
