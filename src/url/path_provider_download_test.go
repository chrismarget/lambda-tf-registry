package url

import (
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/chrismarget/lambda-tf-registry/src/env"
	"github.com/stretchr/testify/require"
)

func TestProviderDownloadPath(t *testing.T) {
	err := os.Setenv("PROVIDER_TABLE_NAME", "registry-providers")
	require.NoError(t, err)
	url := "/v1/providers/hashicorp/tls/4.0.1/download/linux/amd64"

	path, err := NewPathFromString(url)
	require.NoError(t, err)

	sess, err := session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable})
	require.NoError(t, err)

	ddb := dynamodb.New(sess)
	env := env.ParseEnv()

	r, err := path.HandleRequest(ddb, env)
	require.NoError(t, err)

	log.Println(r)
}
