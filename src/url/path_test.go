package url

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderPath_fromString(t *testing.T) {
	type testCase struct {
		path        string
		expPathType PathType
	}

	testCases := map[string]testCase{
		"valid_provider_versions": {
			path:        "/v1/providers/hashicorp/random/versions",
			expPathType: PathTypeProvider,
		},
		"valid_provider_download": {
			path:        "/v1/providers/hashicorp/random/2.0.0/download/linux/amd64",
			expPathType: PathTypeProvider,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			path, err := NewPathFromString(tCase.path)
			require.NoError(t, err)
			assert.Equal(t, tCase.expPathType, path.Type())
		})
	}
}
