package main

import (
	"fmt"
	"os"
)

const (
	namespace       = "jtaf"
	envProviderType = "PTYPE"
)

func getNamespaceType() (string, error) {
	ptype := os.Getenv(envProviderType)
	if ptype == "" {
		return "", fmt.Errorf("environment variable %q must be set", envProviderType)
	}

	return fmt.Sprintf("%s/%s", namespace, ptype), nil
}
