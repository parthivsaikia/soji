package discovery

type Discoverer interface {
	Discover() error
}
