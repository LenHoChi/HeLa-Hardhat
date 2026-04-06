package bank

type impl struct{}

func New() Service {
	return &impl{}
}
