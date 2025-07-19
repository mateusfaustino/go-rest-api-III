package database

import (
	"testing"

	"github.com/mateusfaustino/go-rest-api-III/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserDB(t *testing.T) (*gorm.DB, *UserDb, *entity.Role) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("could not open db: %v", err)
	}
	err = db.AutoMigrate(&entity.Role{}, &entity.User{})
	if err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	role, _ := entity.NewRole("tester")
	db.Create(role)
	return db, NewUserDb(db), role
}

func TestCreateUser(t *testing.T) {
	db, userDB, role := setupUserDB(t)

	user, err := entity.NewUser("Mateus", "m@gmail.com", "123456789", role.ID)
	assert.NoError(t, err)

	err = userDB.CreateUser(user)
	assert.NoError(t, err)

	var userFound entity.User
	err = db.First(&userFound, "id=?", user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)
	assert.NotNil(t, userFound.Password)
}

func TestFindByEmail(t *testing.T) {
	_, userDB, role := setupUserDB(t)

	user, _ := entity.NewUser("Mateus", "m@gmail.com", "123456789", role.ID)
	assert.NoError(t, userDB.CreateUser(user))

	userFound, err := userDB.FindUserByEmail(user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)
	assert.NotNil(t, userFound.Password)
}

func TestFindUserById(t *testing.T) {
	_, userDB, role := setupUserDB(t)

	user, err := entity.NewUser("Mateus", "m@gmail.com", "123456789", role.ID)
	assert.NoError(t, err)
	assert.NoError(t, userDB.CreateUser(user))

	userFound, err := userDB.FindUserById(user.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)
}

func TestUpdateUser(t *testing.T) {
	db, userDB, role1 := setupUserDB(t)
	role2, _ := entity.NewRole("admin")
	db.Create(role2)

	user, err := entity.NewUser("Mateus", "m@gmail.com", "123456789", role1.ID)
	assert.NoError(t, err)
	assert.NoError(t, userDB.CreateUser(user))

	user.Name = "Lucena"
	user.RoleID = role2.ID
	err = userDB.UpdateUser(user)
	assert.NoError(t, err)

	productFound, err := userDB.FindUserById(user.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, "Lucena", productFound.Name)
	assert.Equal(t, role2.ID, productFound.RoleID)
}
