package v2

import (
	"reflect"
	"strings"

	"github.com/meroxa/turbine-go"
)

func (t *Turbine) GetFunction(name string) (turbine.Function, bool) {
	f, ok := t.functions[name]
	return f, ok
}

func (t *Turbine) ListFunctions() []string {
	var funcNames []string
	for name := range t.functions {
		funcNames = append(funcNames, name)
	}
	return funcNames
}

func (t *Turbine) Process(rr turbine.Records, fn turbine.Function) turbine.Records {
	funcName := strings.ToLower(reflect.TypeOf(fn).Name())
	t.functions[funcName] = fn
	t.deploySpec.Functions = append(t.deploySpec.Functions,
		specFunction{Name: funcName, Image: t.imageName})
	return rr
}
