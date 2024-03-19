package responders

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/chrismarget/lambda-tf-registry/src/env"
	"log"
	"net/http"
)

var _ Responder = new(Error)

type Error struct {
	Code int
	Err  error
	Req  *events.LambdaFunctionURLRequest
}

func (o Error) Respond() (events.LambdaFunctionURLResponse, error) {
	response := events.LambdaFunctionURLResponse{StatusCode: o.Code}

	if env.GetBool(env.Debug) {
		var reqString string
		if o.Req != nil {
			b, _ := json.MarshalIndent(o.Req, "", "  ")
			reqString = string(b)
		}

		var errString string
		if o.Err != nil {
			errString = o.Err.Error()
		}

		response.Body = fmt.Sprintf("%d - %q\n\n%s", o.Code, errString, reqString)
	}

	if response.Body != "" {
		log.Print(response.Body)
	}

	var err error
	if o.Code == http.StatusInternalServerError {
		err = o.Err
	}

	return response, err
}
