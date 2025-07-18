package database

import (
	"testing"

	"github.com/mateusfaustino/go-rest-api-III/internal/entity"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateRole(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Role{})

	role, err := entity.NewRole("manager")

	assert.NoError(t, err)

	roleDB := NewRoleDB(db)

	err = roleDB.CreateRole(role)

	assert.NoError(t, err)

	var roleFound entity.Role

	assert.NotEmpty(t, role.ID)

	err = db.First(&roleFound, "id=?", role.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, role.Name, roleFound.Name)
	assert.Equal(t, role.ID, roleFound.ID)

}

func TestFindRoleByEmail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Role{})

	role, _ := entity.NewRole("menager")
	roleDB := NewRoleDB(db)

	err = roleDB.CreateRole(role)

	assert.Nil(t, err)

	roleFound, err := roleDB.FindRoleByName(role.Name)

	assert.Nil(t, err)
	assert.Equal(t, role.ID, roleFound.ID)
	assert.Equal(t, role.Name, roleFound.Name)
}

func TestFindRoleByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Role{})

	role, _ := entity.NewRole("tester")
	roleDB := NewRoleDB(db)
	err = roleDB.CreateRole(role)
	assert.Nil(t, err)

	found, err := roleDB.FindRoleByID(role.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, role.ID, found.ID)
	assert.Equal(t, role.Name, found.Name)
}
