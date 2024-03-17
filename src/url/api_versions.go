package url

import "github.com/orsinium-labs/enum"

type apiVersion enum.Member[string]

var (
	v1          = apiVersion{Value: "v1"}
	apiVersions = enum.New(v1)
)
