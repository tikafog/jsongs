package jsongs

import (
	"fmt"
	"reflect"
)

func methodSet(subv reflect.Value, method reflect.Value) error {
	if method.Kind() != reflect.Func {
		return fmt.Errorf("method %v is not a function", method.Kind())
	}
	method.Call([]reflect.Value{subv})
	return nil
}
