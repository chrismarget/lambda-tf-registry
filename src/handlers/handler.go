package handlers

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
	"regexp"
)

const (
	pathMatchWellKnown = "^/.well-known/"
	pathMatchV1        = "^/v1/"
	pathMatchRegister  = "^/register"
)

type Handler interface {
	Handle(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error)
}

func NewHandlerFromPath(path string) Handler {
	log.Println("enter NewHandlerFromPath")
	switch {
	case regexp.MustCompile(pathMatchWellKnown).MatchString(path):
		log.Println("path matched: ", pathMatchWellKnown)
		return newWellKnownHandler()
	case regexp.MustCompile(pathMatchV1).MatchString(path):
		log.Println("path matched: ", pathMatchV1)
		return newV1Handler()
		//case regexp.MustCompile(pathMatchRegister).MatchString(path):
		//	log.Println("path matched: ", pathMatchRegister)
		//	return newRegisterHandler()
	}

	log.Println("path matched: <none>")
	return newErrorHandler(http.StatusNotFound, errors.New(path))
}
