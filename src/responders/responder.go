package responders

import (
	"github.com/aws/aws-lambda-go/events"
)

type Responder interface {
	Respond() (events.LambdaFunctionURLResponse, error)
}
