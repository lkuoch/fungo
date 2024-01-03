package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

func NewEnvironment() *Environment {
	store := make(map[string]Object)

	return &Environment{
		store: store,
		outer: nil,
	}
}

func (e Environment) Get(name string) (Object, bool) {
	value, ok := e.store[name]
	if !ok && e.outer != nil {
		value, ok = e.outer.Get(name)
	}

	return value, ok
}

func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value

	return value
}
