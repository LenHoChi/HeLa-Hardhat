package blockchain

type gateway struct{}

func New() Gateway {
	return gateway{}
}
