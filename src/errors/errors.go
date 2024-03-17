package ierrors

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

type IErr struct {
	Err  error
	Code int
}

func (o IErr) Error() string {
	return fmt.Sprintf("%d - %s", o.Code, o.Err)
}

func (o IErr) LambdaResponse() (events.LambdaFunctionURLResponse, error) {
	if os.Getenv("DEBUG") == "1" {
		return events.LambdaFunctionURLResponse{
			StatusCode: o.Code,
			Body:       o.Error(),
		}, nil
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: o.Code,
		Body:       strconv.Itoa(o.Code),
	}, nil
}

//func Respond404(err error) (events.LambdaFunctionURLResponse, error) {
//	return events.LambdaFunctionURLResponse{
//		StatusCode: http.StatusNotFound,
//		Body:       err.Error(),
//	}, nil
//}
//
//func Respond500(err error) (events.LambdaFunctionURLResponse, error) {
//	return events.LambdaFunctionURLResponse{
//		StatusCode: http.StatusInternalServerError,
//		Body:       err.Error(),
//	}, nil
//}
