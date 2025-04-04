package entity

import (
	"github.com/mateusfaustino/go-rest-api-III/pkg/entity"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       entity.ID `json:"id" gorm:"type:char(36);primaryKey"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`

	RoleID entity.ID `json:"role_id" gorm:"type:char(36);index"`
	Role   Role      `json:"role" gorm:"foreignKey:RoleID"`
}

func NewUser(name, email, password string, roleID entity.ID) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       entity.NewID(),
		Name:     name,
		Email:    email,
		Password: string(hash),
		RoleID:   roleID,
	}, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
