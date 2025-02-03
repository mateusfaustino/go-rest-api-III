package entity

import (
	"github.com/mateusfaustino/go-rest-api-III/pkg/entity"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       entity.ID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	Password string    `json:"-"`
}

func NewUser(name, email, password string, role string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil { // Handle errors reading the config file
		return nil, err
	}

	// Definir a role padrão como "customer", caso não seja informada
	if role == "" {
		role = "customer"
	}

	return &User{
		ID:       entity.NewID(),
		Name:     name,
		Email:    email,
		Password: string(hash),
		Role:     role,
	}, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
