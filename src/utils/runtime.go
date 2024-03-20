package utils

import "github.com/chrismarget/lambda-tf-registry/src/env"

func runningInAws() bool {
	return env.Get(env.LambdaVar) != ""
}
