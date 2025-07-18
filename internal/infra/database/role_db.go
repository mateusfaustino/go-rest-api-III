package database

import (
	"github.com/mateusfaustino/go-rest-api-III/internal/entity"
	"gorm.io/gorm"
)

type RoleDB struct {
	DB *gorm.DB
}

func NewRoleDB(db *gorm.DB) *RoleDB {
	return &RoleDB{
		DB: db,
	}
}

func (rdb *RoleDB) CreateRole(role *entity.Role) error {
	return rdb.DB.Create(role).Error
}

func (rdb *RoleDB) FindRoleByName(name string) (*entity.Role, error) {
	var role entity.Role
	err := rdb.DB.Where("name = ? ", name).First(&role).Error

	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (rdb *RoleDB) RoleExists(name string) (bool, error) {
        var count int64
        err := rdb.DB.Model(&entity.Role{}).Where("name = ?", name).Count(&count).Error
        if err != nil {
                return false, err
        }
        return count > 0, nil
}

func (rdb *RoleDB) FindRoleByID(id string) (*entity.Role, error) {
        var role entity.Role
        err := rdb.DB.First(&role, "id = ?", id).Error
        if err != nil {
                return nil, err
        }
        return &role, nil
}
