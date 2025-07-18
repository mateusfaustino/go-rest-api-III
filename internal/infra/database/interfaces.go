package database

import "github.com/mateusfaustino/go-rest-api-III/internal/entity"

type UserInterface interface {
	CreateUser(user *entity.User) error
	FindUserByEmail(email string) (*entity.User, error)
	FindUserById(id string) (*entity.User, error)
	UpdateUser(user *entity.User) error
}

type ProductInterface interface {
	CreateProduct(product *entity.Product) error
	FindAllProducts(page, limit int, sort string) ([]entity.Product, error)
	FindProductByID(id string) (*entity.Product, error)
	UpdateProduct(product *entity.Product) error
	DeleteProduct(id string) error
}

type RoleInterface interface {
       FindRoleByName(name string) (*entity.Role, error)
       CreateRole(role *entity.Role) error
       RoleExists(roleName string) (bool, error)
       FindRoleByID(id string) (*entity.Role, error)
}
