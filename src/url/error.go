package url

import "fmt"

type pathError struct {
	responseCode int
	errMsg       string
}

func (o pathError) Error() string {
	return o.errMsg
}

func newPathError(rc int, format string, args ...any) pathError {
	return pathError{
		responseCode: rc,
		errMsg:       fmt.Sprintf(format, args...),
	}
}
