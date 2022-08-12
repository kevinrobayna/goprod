package products

type Product struct {
	Id          uint   `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
}

func NewService(repository IRepository) IService {
	return &Service{repository: repository}
}

type IService interface {
	SaveProduct(product Product) error
	GetProducts() ([]Product, error)
	GetProduct(id uint) (Product, error)
}

type Service struct {
	repository IRepository
}

func (s *Service) SaveProduct(product Product) error {
	domain := mapModelToDomain(product)
	if product.Id != 0 {
		domain.ID = product.Id
	}
	return s.repository.SaveProduct(domain)
}

func (s *Service) GetProducts() ([]Product, error) {
	domains, err := s.repository.GetProducts()
	if err != nil {
		return nil, err
	}
	products := make([]Product, len(domains))
	for i, domain := range domains {
		products[i] = mapDomainToModel(domain)
	}
	return products, nil
}

func (s *Service) GetProduct(id uint) (Product, error) {
	domain, err := s.repository.GetProduct(id)
	if err != nil {
		return Product{}, err
	}
	return mapDomainToModel(domain), nil
}

func mapModelToDomain(product Product) ProductDomain {
	return ProductDomain{
		Code:        product.Code,
		Description: product.Description,
		Price:       product.Price,
	}
}

func mapDomainToModel(product ProductDomain) Product {
	return Product{
		Id:          product.ID,
		Code:        product.Code,
		Description: product.Description,
		Price:       product.Price,
	}
}
