package responders

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

var _ Responder = new(WellKnownTerraform)

type WellKnownTerraform struct {
	ModulesV1   string `json:"modules.v1,omitempty"`
	ProvidersV1 string `json:"providers.v1,omitempty"`
}

func (o WellKnownTerraform) Respond() (events.LambdaFunctionURLResponse, error) {
	body, err := json.Marshal(o)
	if err != nil {
		return Error{
			Code: http.StatusInternalServerError,
			Err:  err,
		}.Respond()
	}

	response := events.LambdaFunctionURLResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}

	return response, nil
}

//func NewWellKnownResponder(v1Modules, v1Providers string) Responder {
//	return WellKnownTerraform{
//		ModulesV1:   v1Modules,
//		ProvidersV1: v1Providers,
//	}
//}
