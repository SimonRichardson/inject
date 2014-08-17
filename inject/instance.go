package inject

import "reflect"

func instanceOf(t Any) Any {
	x := reflect.ValueOf(t)
	return x.Interface()
}

func isInstanceOf(module IModule, t Any) bool {
	if module == nil || t == nil {
		return module == t
	}
	v1 := reflect.ValueOf(module)
	v2 := reflect.ValueOf(t)
	if v1.Type() != v2.Type() {
		return false
	}
	return false
}

func instanceEquals(x Any, y Any) bool {
	return x == y
}
