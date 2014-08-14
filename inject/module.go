package inject

type Module struct {
	bindings    map[Any]*Binding
	initialised bool
}

func NewModule() *Module {
	return &Module{
		bindings:    make(map[Any]*Binding),
		initialised: false,
	}
}

func (m Module) Initialise() {
	m.Configure()
	m.initialised = true
}

func (m Module) Configure() {}

func (m Module) GetInstance(t Any) Option {
	if !m.initialised {
		return NewNone()
	}

	PushScope(m)
	defer func() {
		PopScope()
	}()

	return find(t).Fold(
		func(x Any) Option {
			return NewSome(x.GetInstance())
		},
		func() {
			return NewNone()
		},
	)
}

func (m Module) Binds(t Any) bool {
	return m.Find(t).Bool()
}

func (m Module) Bind(t Any) *Binding {
	return m.BindWith(t, make([]Any, 0, 0))
}

func (m Module) BindWith(t Any, a []Any) *Binding {

	binding := NewBinding(m)
	binding.To(t)

	bindings[t] = binding

	return binding
}

func (m Module) Unbind(t Any) {
	delete(m, t)
}

func (m Module) Find(t Any) Option {
	if x, ok := m.bindings[t]; ok {
		return NewSome(x)
	}
	return NewNone()
}

func (m Module) Dispose() {
	bindings = make(map[Any]*Binding)
}
