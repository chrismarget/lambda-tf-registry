package handlers

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
	"strings"
)

var _ Handler = new(V1Handler)

const (
	v1PathProvider = "providers"
	v1PathModule   = "modules"
)

type V1Handler struct{}

// /v1/modules/hashicorp/consul/aws/versions
// /v1/modules/hashicorp/consul/aws/0.0.1/download

func (v V1Handler) Handle(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	log.Println("enter V1Handler.Handle()")
	urlParts := strings.Split(strings.TrimLeft(req.RawPath, "/"), "/")
	log.Printf("part 1 is: %q", urlParts[1])

	var h Handler
	switch {
	case len(urlParts) < 2:
		log.Println("insufficient URL parts part count")
		h = newErrorHandler(http.StatusNotFound, fmt.Errorf("insufficient URL parts (%d) for v1 handler", len(urlParts)))
	//case urlParts[1] == v1PathModule: // todo
	//	log.Println("create V1ModuleHandler")
	//	h = newV1ModuleHandler(urlParts)
	case urlParts[1] == v1PathProvider:
		log.Println("create V1ProviderHandler")
		h = newV1ProviderHandler(urlParts)
	default:
		log.Printf("unknown v1 path %q\n", req.RawPath)
		h = newErrorHandler(http.StatusNotFound, nil)
	}

	return h.Handle(ctx, req)
}

func newV1Handler() Handler {
	return V1Handler{}
}

//func newV1Handler(urlParts []string) Handler {
//	log.Println("enter newV1Handler")
//	switch {
//	case urlParts[1] == v1PathProvider:
//		log.Printf("creating %q handler", v1PathProvider)
//		return newV1ProviderHandler(urlParts)
//	case urlParts[1] == v1PathModule:
//		log.Printf("creating %q handler", v1PathModule)
//		return newErrorHandler(http.StatusNotFound, fmt.Errorf("v1 path %q not handled", v1PathModule))
//	}
//
//	return newErrorHandler(http.StatusNotFound, fmt.Errorf("v1 path %q not handled", urlParts[1]))
//}
