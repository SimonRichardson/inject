package inject

var (
	currentScope Option          = NewNone()
	modules      []IModule       = make([]IModule, 0, 0)
	scopes       []IModule       = make([]IModule, 0, 0)
	bindings     map[Any]IModule = make(map[Any]IModule)
)

func Add(module IModule) IModule {
	modules = append(modules, module)
	module.Initialise()
	return module
}

func Remove(module IModule) IModule {
	module.Dispose()

	modules = remove(modules, module)
	scopes = remove(scopes, module)

	currentScope = currentScope.Chain(func(x Any) Option {
		if instanceEquals(x, module) {
			return NewNone()
		}
		return NewSome(x)
	})

	return module
}

func PushScope(module IModule) {
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
	binding := bindings[t]
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

func remove(x []IModule, mod IModule) []IModule {
	index := -1
	for k, v := range x {
		if instanceEquals(v, mod) {
			index = k
			break
		}
	}

	if index == -1 {
		return x
	}

	return append(x[:index], x[index+1:]...)
}

func removeLast(x []IModule) []IModule {
	index := len(x)
	if index < 1 {
		return x
	}
	return append(x[:index], x[index+1:]...)
}
