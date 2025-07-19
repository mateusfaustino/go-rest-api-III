package database

import (
	"github.com/mateusfaustino/go-rest-api-III/internal/entity"
	"gorm.io/gorm"
)

type UserDb struct {
	DB *gorm.DB
}

func NewUserDb(db *gorm.DB) *UserDb {
	return &UserDb{
		DB: db,
	}
}

func (u *UserDb) CreateUser(user *entity.User) error {
	return u.DB.Create(user).Error
}

func (u *UserDb) FindUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := u.DB.Where("email = ? ", email).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserDb) FindUserById(id string) (*entity.User, error) {
	var user entity.User
	err := u.DB.First(&user, "id=?", id).Error
	return &user, err
}

func (u *UserDb) UpdateUser(user *entity.User) error {
	return u.DB.Save(user).Error
}
