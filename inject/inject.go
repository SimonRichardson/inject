package inject

import "fmt"

func As(t Any) Option {
	fmt.Println(">>>", CurrentScope())
	return CurrentScope().Fold(
		func(x Any) Option {
			instance := x.(IModule).GetInstance(t)
			return instance.Fold(
				func(x Any) Option {
					return instance
				},
				func() Option {
					return With(t)
				},
			)
		},
		func() Option {
			return With(t)
		},
	)
}

func With(t Any) Option {
	return ScopeOf(t).Chain(func(x Any) Option {
		return x.(Module).GetInstance(t)
	})
}

func WithIn(t Any, module Any) Option {
	return ModuleOf(module).Chain(func(x Any) Option {
		return x.(Module).GetInstance(t)
	})
}
