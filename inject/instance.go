package inject

import "reflect"

func instanceOf(t Any) Any {
	x := reflect.ValueOf(x)
	return x.Interface()
}

func isInstanceOf(module *Module, t Any) bool {
	if _, ok := module.(t); ok {
		return true
	}
	return false
}
