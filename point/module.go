package point

import "fmt"

const (
	ErrBindingErrorInit = fmt.Errorf("Modules have to be created using \"Add(NewModule())\".")
	ErrInstanceMiss     = fmt.Errorf("Type not found.")
)

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

func (m Module) Initialised() {
	m.Configure()
	m.initialised = true
}

func (m Module) Configure() {}

func (m Module) GetInstance(t Any) Either {
	if !m.initialised {
		return NewLeft(ErrBindingErrorInit)
	}

	PushScope(m)

	instance := find(t).Fold(
		func(x Binding) {
			return NewRight(x.GetInstance())
		},
		func() {
			return NewLeft(ErrInstanceMiss)
		},
	)

	PopScope()

	return instance
}

func (m Module) Binds(t Any) bool {
	return find(t).Bool()
}

func (m Module) Bind(t Any) *Binding {
	return m.BindWith(t, make([]Any, 0, 0))
}

func (m Module) BindWith(t Any, a []Any) *Binding {
	return nil
}

func (m Module) Unbind(t Any) {

}

func find(t Any) Option {
	return None
}
