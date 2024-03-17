package url

import (
	"github.com/orsinium-labs/enum"
)

type PathType enum.Member[string]

func (o PathType) String() string { return o.Value }

var (
	PathTypeModule   = PathType{Value: "modules"}
	PathTypeProvider = PathType{Value: "providers"}
	PathTypes        = enum.New(
		PathTypeModule,
		PathTypeProvider,
	)
)
