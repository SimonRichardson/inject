package inject

type IModule interface {
	Initialise()
	GetInstance(Any) Option
	Binds(Any) bool
	Bind(Any) *Binding
	BindWith(Any, []Any) *Binding
	Unbind(Any)
	Find(Any) Option
	Dispose()
}

type Module struct {
	bindings    map[Any]*Binding
	initialised bool
	configure   func(IModule)
}

func NewModule(configure func(IModule)) *Module {
	return &Module{
		bindings:    make(map[Any]*Binding),
		initialised: false,
		configure:   configure,
	}
}

func (m *Module) Initialise() {
	m.configure(m)
	m.initialised = true
}

func (m *Module) GetInstance(t Any) Option {
	if !m.initialised {
		return NewNone()
	}

	PushScope(m)
	defer func() {
		PopScope()
	}()

	return m.Find(t).Fold(
		func(x Any) Option {
			return NewSome(x.(*Binding).GetInstance())
		},
		func() Option {
			if c, ok := t.(Binder); ok {
				return NewSome(c.Get())
			}
			return NewNone()
		},
	)
}

func (m *Module) Binds(t Any) bool {
	return m.Find(t).Bool()
}

func (m *Module) Bind(t Any) *Binding {
	return m.BindWith(t, make([]Any, 0, 0))
}

func (m *Module) BindWith(t Any, a []Any) *Binding {
	binding := NewBinding(m)
	binding.To(t)

	m.bindings[t] = binding

	return binding
}

func (m *Module) Unbind(t Any) {
	delete(m.bindings, t)
}

func (m *Module) Find(t Any) Option {
	if x, ok := m.bindings[t]; ok {
		return NewSome(x)
	}
	return NewNone()
}

func (m *Module) Dispose() {
	m.bindings = make(map[Any]*Binding)
}
