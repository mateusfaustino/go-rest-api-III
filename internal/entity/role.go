package entity

import "github.com/mateusfaustino/go-rest-api-III/pkg/entity"

type Role struct {
	ID    entity.ID `json:"id" gorm:"type:char(36);primaryKey"`
	Name  string    `json:"name" gorm:"unique;not null"`
	Users []User    `json:"users" gorm:"foreignKey:RoleID"`
}

func NewRole(name string) (*Role, error) {

	return &Role{
		ID:   entity.NewID(),
		Name: name,
	}, nil
}
