package utils

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/chrismarget/lambda-tf-registry/src/env"
	"net/http"
	"os"
)

var unsafeAuthTestModeOK bool // set only in integration tests

func authTestModeEnabled() bool {
	return unsafeAuthTestModeOK
}

const (
	AuthHeader = "AuthToken"
)

func SecretsManagerToken() (string, error) {
	secretId := env.Get(env.RegisterTokenName)
	if secretId == "" {
		return "", HttpResponseErr{
			err:    errors.New("cannot authenticate - unknown secret ID"),
			status: http.StatusInternalServerError,
		}
	}

	client, err := SmClient()
	if err != nil {
		return "", HttpResponseErr{
			err:    fmt.Errorf("cannot authenticate - while getting secrets manager client - %w", err),
			status: http.StatusInternalServerError,
		}
	}

	gsvo, err := client.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: &secretId})
	if err != nil {
		return "", HttpResponseErr{
			err:    fmt.Errorf("cannot authenticate - while getting secret value - %w", err),
			status: http.StatusInternalServerError,
		}
	}
	if gsvo.SecretString == nil || *gsvo.SecretString == "" {
		return "", HttpResponseErr{
			err:    errors.New("cannot authenticate - secret value not set"),
			status: http.StatusInternalServerError,
		}
	}

	return *gsvo.SecretString, nil
}

func testToken() (string, error) {
	token := env.Get(AuthHeader)
	if token == "" {
		return "", HttpResponseErr{
			err:    errors.New("cannot authenticate - failed loading auth token"),
			status: http.StatusUnauthorized,
		}
	}

	return token, nil
}

func testAuthTokenSet() bool {
	_, ok := os.LookupEnv(AuthHeader)
	return ok
}

func authTestEnabled() bool {
	return !runningInAws() && authTestModeEnabled() && testAuthTokenSet()
}

func ValidateToken(userToken string) error {
	if userToken == "" {
		return HttpResponseErr{
			err:    fmt.Errorf("%s http header required", AuthHeader),
			status: http.StatusUnauthorized,
		}
	}

	var authToken string
	var err error
	switch {
	case runningInAws():
		authToken, err = SecretsManagerToken()
	case authTestEnabled():
		authToken, err = testToken()
	}
	if err != nil {
		return err
	}

	if userToken == authToken {
		return nil
	}

	return HttpResponseErr{
		err:    errors.New("authorization failed"),
		status: http.StatusUnauthorized,
	}
}
