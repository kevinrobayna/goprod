package products

import (
	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{db: db}
}

func InvokeMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&ProductDomain{})
	if err != nil {
		return
	}
}

type ProductDomain struct {
	gorm.Model
	Code        string
	Description string
	Price       uint
}

type IRepository interface {
	SaveProduct(product ProductDomain) error
	GetProducts() ([]ProductDomain, error)
	GetProduct(id uint) (ProductDomain, error)
}

type Repository struct {
	db *gorm.DB
}

func (r *Repository) SaveProduct(product ProductDomain) error {
	return r.db.Save(&product).Error
}

func (r *Repository) GetProducts() ([]ProductDomain, error) {
	var products []ProductDomain
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *Repository) GetProduct(id uint) (ProductDomain, error) {
	var product ProductDomain
	if err := r.db.Take(&product, id).Error; err != nil {
		return ProductDomain{}, err
	}
	return product, nil
}
