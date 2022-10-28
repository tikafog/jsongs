package jsongs

import (
	"reflect"
	"unicode"
)

func methodName(name string, tag string, prefix string) string {
	if tag != "" {
		return tag
	}
	name = camelcaseMethodName(name)
	if prefix != "" {
		name = prefix + name
	}
	return name
}

func methodFuncValue(fv reflect.Value, v reflect.Value, name string) reflect.Value {
	if name == "" {
		return fv
	}
	method := v.MethodByName(name)
	if method.Kind() != reflect.Func && v.Kind() != reflect.Pointer && v.CanAddr() {
		method = v.Addr().MethodByName(name)
	} else if method.Kind() != reflect.Func && v.Kind() != reflect.Pointer && !v.CanAddr() {
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

func methodFuncName(v reflect.Type, name string) string {
	method, exist := methodFuncMethod(v, name)
	if exist {
		return method.Name
	}
	return ""
}

func camelcaseMethodName(name string) string {
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
