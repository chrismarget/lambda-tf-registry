package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var _ error = new(handlerError)

type handlerError struct {
	statusCode int
	errPublic  error
	errPrivate error
	headers    map[string]string
}

func (o *handlerError) Error() string {
	return o.errPublic.Error()
}

func (o *handlerError) SetCode(code int) {
	o.statusCode = code
}

func (o *handlerError) SetHeaders(headers map[string]string) {
	o.headers = headers
}

func (o handlerError) MarshalResponse() (events.APIGatewayProxyResponse, error) {
	message := json.RawMessage(o.errPublic.Error())
	if GetBool(Debug) {
		message = json.RawMessage(fmt.Sprintf(`{"public":%q,"private":%q}`, o.errPublic, o.errPrivate))
	}

	data := struct {
		Code    int             `json:"code"`
		Message json.RawMessage `json:"message"`
	}{
		Code:    o.statusCode,
		Message: message,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    map[string]string{"Content-Type": "application/json; charset=UTF-8"},
			Body:       `{"code":500,"message":"the server has encountered an error while handling a different error"}`,
		}, err
	}

	if o.headers == nil {
		o.headers = map[string]string{"Content-Type": "application/json; charset=UTF-8"}
	} else {
		o.headers["Content-Type"] = "application/json; charset=UTF-8"
	}

	return events.APIGatewayProxyResponse{
		StatusCode: o.statusCode,
		Headers:    o.headers,
		Body:       string(body),
	}, nil
}

func FromPrivateError(private error, pubMsg string) handlerError {
	return handlerError{
		statusCode: http.StatusInternalServerError,
		errPublic:  errors.New(pubMsg),
		errPrivate: private,
		headers:    nil,
	}
}

func FromPublicError(err error) handlerError {
	return handlerError{
		statusCode: http.StatusInternalServerError,
		errPublic:  err,
		errPrivate: err,
		headers:    nil,
	}
}

//func FromPrivateError(private error, public ...error) handlerError {
//	var pubErr error
//	if len(public) == 0 {
//		pubErr = errors.New("see logs for details")
//	} else {
//		pubErr = errors.Join(public...)
//	}
//	return handlerError{
//		statusCode: http.StatusInternalServerError,
//		errPublic:  pubErr,
//		errPrivate: private,
//		headers:    nil,
//	}
//}

func NewPublicError(s string) handlerError {
	return FromPublicError(errors.New(s))
}
