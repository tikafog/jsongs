package jsongs

import (
	"reflect"
	"unicode"
)

func methodName(t reflect.Type, name string, tag string, prefix string) string {
	if tag != "" {
		return tag
	}
	name = methodNameFix(name)
	if prefix != "" {
		name = prefix + name
	}
	return name
}

func checkMethodExists(t reflect.Type, name string) bool {
	_, ok := t.MethodByName(name)
	if !ok {
		_, ok = reflect.PtrTo(t).MethodByName(name)
	}
	return ok
}

func methodNameIfExists(t reflect.Type, name string) (string, bool) {
	_, ok := t.MethodByName(name)
	if !ok {
		_, ok = reflect.PtrTo(t).MethodByName(name)
	}
	return name, ok
}

func methodFuncValue(v reflect.Value, name string) reflect.Value {
	var method reflect.Value
	method = v.MethodByName(name)
	if method.Kind() != reflect.Func && v.CanAddr() {
		method = v.Addr().MethodByName(name)
	} else if method.Kind() != reflect.Func && !v.CanAddr() && v.Kind() != reflect.Pointer {
		nv := reflect.New(v.Type())
		nv.Elem().Set(v)
		method = nv.MethodByName(name)
	} else {

	}

	return method
}

func methodFuncMethod(v reflect.Type, name string) (reflect.Method, bool) {
	var method reflect.Method
	var exist bool
	method, exist = v.MethodByName(name)
	if !exist {
		method, exist = reflect.PtrTo(v).MethodByName(name)
	}
	return method, exist
}

func methodNameFix(name string) string {
	nameRunes := []rune(name)
	var ret []rune
	if unicode.IsLower(nameRunes[0]) {
		ret = append(ret, unicode.ToUpper(nameRunes[0]))
		nameRunes = nameRunes[1:]
	}
	nextUpper := false
	for i := range nameRunes {
		if nameRunes[i] == rune('_') {
			nextUpper = true
			continue
		}
		if nextUpper {
			nextUpper = false
			ret = append(ret, unicode.ToUpper(nameRunes[i]))
			continue
		}
		ret = append(ret, nameRunes[i])
	}
	return string(ret)
}
