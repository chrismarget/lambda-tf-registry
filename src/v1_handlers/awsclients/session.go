package awsclients

import "github.com/aws/aws-sdk-go/aws/session"

var sess *session.Session

func getSession() (*session.Session, error) {
	var err error
	if sess == nil {
		sess, err = session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable})
	}

	return sess, err
}
