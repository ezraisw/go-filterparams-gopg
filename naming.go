package fpgopg

type Naming interface {
	Interpret(name string) string
}

type NamingFunc func(string) string

func (n NamingFunc) Interpret(name string) string {
	return n(name)
}

func PassthroughNaming(name string) string {
	return name
}
