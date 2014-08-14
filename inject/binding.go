package inject

type BindingType int

const (
	bindingTypeTo       BindingType = 1
	bindingTypeInstance BindingType = 2
	bindingTypeProvider BindingType = 3
)

type Provider interface {
	InstanceOf() Any
}

type Binder interface {
	Get() Any
}

type to struct {
	Type Any
}

func (b to) Get() Any {
	return b.Type
}

type instance struct {
	Instance Any
}

func (b instance) Get() Any {
	return b.Instance
}

type provider struct {
	Provider Provider
}

func (b provider) Get() Any {
	return b.Provider.InstanceOf()
}

type Binding struct {
	module         *Module
	bindingType    BindingType
	binding        *Binder
	singletonScope bool
	evaluated      bool
	value          Any
}

func NewBinding(module *Module) *Binding {
	return &Binding{
		module:         module,
		evaluated:      false,
		singletonScope: false,
	}
}

func (b *Binding) To(t Any) Scope {
	b.evaluated = false
	b.bindingType = bindingTypeTo
	b.binding = &to{Type: t}

	return b
}

func (b *Binding) ToInstance(i Any) Scope {
	b.evaluated = false
	b.bindingType = bindingTypeInstance
	b.binding = &instance{Instance: i}

	return b
}

func (b *Binding) ToProvider(p Provider) Scope {
	b.evaluated = false
	b.bindingType = bindingTypeProvider
	b.binding = &provider{Provider: p}

	return b
}

func (b *Binding) GetInstance() Any {
	if b.singletonScope {
		if b.evaluated {
			b.value = solve()
			b.evaluated = true
		}
		return b.value
	}
	return solve()
}

func (b *Binding) AsSingleton() {
	b.singletonScope = true
}

func (b *Binding) solve() Any {
	switch b.bindingType {
	case bindingTypeTo:
		return instanceOf(b.binding.Get())
	case bindingTypeInstance:
		return b.binding.Get()
	case bindingTypeProvider:
		return b.binding.Get()
	}
}
