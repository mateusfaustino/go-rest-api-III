package database

import "github.com/mateusfaustino/go-rest-api-III/internal/entity"

type UserInterface interface {
	CreateUser(user *entity.User) error
	FindUserByEmail(email string) (*entity.User, error)
}