package internal

import (
	"gorm.io/gorm"
)

type product struct {
	gorm.Model
	ID          uint `gorm:"primarykey"`
	Code        string
	Description string
	Price       uint
}

type iRepository interface {
	SaveProduct(product product) (product, error)
	UpdateProduct(product product) (product, error)
	GetProducts() ([]product, error)
	GetProduct(id uint) (product, error)
}

type repository struct {
	db *gorm.DB
}

func (r *repository) SaveProduct(p product) (product, error) {
	if err := r.db.Create(&p).Error; err != nil {
		return product{}, err
	}
	return p, nil
}

func (r *repository) UpdateProduct(p product) (product, error) {
	if err := r.db.Updates(&p).Error; err != nil {
		return product{}, err
	}
	return p, nil
}

func (r *repository) GetProducts() ([]product, error) {
	var products []product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *repository) GetProduct(id uint) (product, error) {
	var p product
	if err := r.db.Take(&p, id).Error; err != nil {
		return product{}, err
	}
	return p, nil
}
