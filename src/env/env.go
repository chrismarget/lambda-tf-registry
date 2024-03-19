package env

import (
	"os"
	"strconv"
)

const (
	Debug             = "DEBUG"
	ProviderTableName = "PROVIDER_TABLE_NAME"
)

var defaults = map[string]string{
	Debug:             "false",
	ProviderTableName: "registry-providers",
}

type Env struct { // todo: make private
	vars     map[string]string
	boolVars map[string]bool
}

func (o *Env) Get(s string) string {
	if r, ok := o.vars[s]; ok {
		return r
	}

	if v, ok := os.LookupEnv(s); ok {
		o.vars[s] = v
		return v
	}

	if v, ok := defaults[s]; ok {
		o.vars[s] = v
		return v
	}

	return ""
}

func (o *Env) GetBool(s string) bool {
	if v, ok := o.boolVars[s]; ok {
		return v
	}

	v, _ := strconv.ParseBool(o.Get(s))
	o.boolVars[s] = v
	return v
}

var env = Env{
	vars:     map[string]string{},
	boolVars: map[string]bool{},
}

func Get(s string) string {
	return env.Get(s)
}

func GetBool(s string) bool {
	return env.GetBool(s)
}
