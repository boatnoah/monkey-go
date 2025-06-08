package object

func NewEnviroment() *Enviroment {
	s := make(map[string]Object)
	return &Enviroment{store: s}
}

type Enviroment struct {
	store map[string]Object
}

func (e *Enviroment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Enviroment) Set(name string, obj Object) Object {

	e.store[name] = obj
	return obj

}
