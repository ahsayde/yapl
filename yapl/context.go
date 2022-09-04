package yapl

import (
	"os"
	"strings"
)

type Context struct {
	Input  interface{}
	Params map[string]interface{}
	Env    map[string]string
	Cond   ConditionResult
}

func newContext(input, params map[string]interface{}) *Context {
	return &Context{
		Input:  input,
		Env:    environ(),
		Params: params,
	}
}

func environ() map[string]string {
	vars := make(map[string]string)
	for _, v := range os.Environ() {
		part := strings.Split(v, "=")
		vars[part[0]] = part[1]
	}
	return vars
}
