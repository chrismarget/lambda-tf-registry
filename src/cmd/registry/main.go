package main

import (
	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chrismarget/lambda-tf-registry/src/v1_handlers"
)

var router = lmdrouter.NewRouter("")

func init() {
	serviceDiscoveryMap := map[string]string{
		"v1.providers": "/v1/providers/",
		//"v1.modules":   "/v1/modules/",
	}

	v1handlers.NewServiceDiscoveryHandler(serviceDiscoveryMap).AddRoutes(router)
	v1handlers.NewProviderDownloadHandler().AddRoutes(router)
	v1handlers.NewProviderVersionsHandler().AddRoutes(router)
}

func main() {
	lambda.Start(router.Handler)
}
