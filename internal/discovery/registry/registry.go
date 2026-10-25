package registry

type Registry interface {
	Discover() ([]string, error)
}

func Discover(reg Registry) ([]string, error) {
	return reg.Discover()
}
