//go:build tools

package tools

import (
	_ "github.com/google/go-licenses"
	_ "honnef.co/go/tools/cmd/staticcheck"
	_ "mvdan.cc/gofumpt"
)
