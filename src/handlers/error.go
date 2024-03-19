package handlers

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/chrismarget/lambda-tf-registry/src/responders"
)

var _ Handler = new(Error)

type Error struct {
	code int
	err  error
}

func (o Error) Handle(_ context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	responder := responders.Error{
		Code: o.code,
		Err:  o.err,
		Req:  &req,
	}

	return responder.Respond()
	//response := events.LambdaFunctionURLResponse{StatusCode: o.code}
	//
	//if env.GetBool(env.Debug) {
	//	h, _ := json.MarshalIndent(req.Headers, "", "  ")
	//	response.Body = fmt.Sprintf(
	//		"\n"+
	//			"%d\n\n"+
	//			"%s\n\n"+
	//			"%s\n\n",
	//		o.code,
	//		o.err.Error(),
	//		string(h),
	//	)
	//}
	//
	//log.Println(fmt.Sprintf("%d - %q", o.code, o.err))
	//
	//return response, nil
}

func newErrorHandler(code int, err error) Handler {
	return Error{
		code: code,
		err:  err,
	}
}
