package database

import (
	"github.com/mateusfaustino/go-rest-api-III/internal/entity"
	"gorm.io/gorm"
)

type User struct {
	DB *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{
		DB: db,
	}
}

func (u *User) CreateUser(user *entity.User) error {
	return u.DB.Create(user).Error
}

func (u *User) FindUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := u.DB.Where("email = ? ", email).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
