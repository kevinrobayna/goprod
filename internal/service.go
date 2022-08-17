package internal

type IService interface {
	SaveProduct(product Product) error
	GetProducts() ([]Product, error)
	GetProduct(id uint) (Product, error)
}

type Service struct {
	repository iRepository
}

func (s Service) SaveProduct(product Product) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetProducts() ([]Product, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetProduct(id uint) (Product, error) {
	//TODO implement me
	panic("implement me")
}
