package utils

type HttpResponseErr struct {
	err    error
	status int
}

func (o HttpResponseErr) Error() string {
	return o.err.Error()
}

func (o HttpResponseErr) HttpStatusCode() int {
	return o.status
}
