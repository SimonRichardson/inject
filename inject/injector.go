package inject

var (
	currentScope Option          = NewNone()
	modules      []*Module       = new([]*Module, 0, 0)
	scopes       []*Module       = new([]*Module, 0, 0)
	bindings     map[Any]*Module = make(map[Any]*Module)
)

func Add(module Module) Module {
	modules = append(modules, module)
	module.Initialise()
	return module
}

func Remove(module Module) Module {
	module.Dispose()

	modules = remove(modules, module)
	scopes = remove(scopes, module)

	currentScope = currentScope.Map(func(x Any) Option {
		if x.(Module) == module {
			return NewNone()
		}
		return x
	})

	return module
}

func PushScope(module Module) {
	currentScope = NewSome(module)
	scopes = append(scopes, module)
}

func PopScope() {
	scopes = removeLast(scopes)
	numOfScopes := len(scopes)
	if numOfScopes == 0 {
		currentScope = NewNone()
	} else {
		currentScope = NewSome(scopes[numOfScopes])
	}
}

func CurrentScope() Option {
	return currentScope
}

func ScopeOf(t Any) Option {
	res := NewNone()

	for _, module := range modules {
		if module.Binds(t) {
			return NewSome(module)
		}
	}

	return res
}

func ModuleOf(t Any) Option {
	binding = bindings[t]
	if binding != nil {
		return NewSome(binding)
	}

	for _, module := range modules {
		if isInstanceOf(module, t) {
			bindings[t] = module
			return NewSome(module)
		}
	}

	return NewNone()
}

func remove(x []*Module, mod Module) []*Module {
	index := -1
	for k, v := range x {
		if v == mod {
			index = k
			break
		}
	}

	if index == -1 {
		return x
	}

	return append(x[:index], s[index+1:]...)
}

func removeLast(x []*Module) []*Module {
	index := len(x)
	if index < 1 {
		return x
	}
	return append(x[:index], s[index+1:]...)
}
