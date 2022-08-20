package internal

import "gorm.io/gorm"

type Product struct {
	Id          uint   `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
}

type IService interface {
	SaveProduct(product Product) error
	GetProducts() ([]Product, error)
	GetProduct(id uint) (Product, error)
}

type Service struct {
	repository iRepository
}

func (s Service) SaveProduct(product Product) error {
	if product.Id != 0 {
		_, err := s.repository.UpdateProduct(mapModelToDao(product))
		if err != nil {
			return err
		}
	} else {
		_, err := s.repository.SaveProduct(mapModelToDao(product))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Service) GetProducts() ([]Product, error) {
	productsDao, err := s.repository.GetProducts()
	if err != nil {
		return []Product{}, err
	}
	var products []Product
	for _, p := range productsDao {
		products = append(products, mapDaoToModel(p))
	}
	return products, nil
}

func (s Service) GetProduct(id uint) (Product, error) {
	p, err := s.repository.GetProduct(id)
	if err != nil {
		return Product{}, err
	}
	return mapDaoToModel(p), nil
}

func mapDaoToModel(dao product) Product {
	return Product{
		Id:          dao.ID,
		Code:        dao.Code,
		Description: dao.Description,
		Price:       dao.Price,
	}
}

func mapModelToDao(model Product) product {
	return product{
		Model:       gorm.Model{ID: model.Id},
		Code:        model.Code,
		Description: model.Description,
		Price:       model.Price,
	}
}
