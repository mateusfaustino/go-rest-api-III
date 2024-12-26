package database

import "github.com/mateusfaustino/go-rest-api-III/internal/entity"

type UserInterface interface {
	CreateUser(user *entity.User) error
	FindUserByEmail(email string) (*entity.User, error)
}

type ProductInterface interface {
	CreateProduct(product *entity.Product) error
	FindAllProducts(page, limit int, sort string) ([]entity.Product, error)
	FindProductByID(id string) (*entity.Product, error)
	UpdateProduct(product *entity.Product) error
	DeleteProduct(id string) error
}
