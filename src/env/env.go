package env

import "os"

const ProviderTableName = "PROVIDER_TABLE_NAME"

type Env struct {
	TableName string
}

func ParseEnv() Env {
	return Env{
		TableName: os.Getenv(ProviderTableName),
	}
}
