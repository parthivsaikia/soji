package registry

type Registry interface {
	Discover() ([]string, error)
}
