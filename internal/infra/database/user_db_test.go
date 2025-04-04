package database

// import (
// 	"testing"

// 	"github.com/mateusfaustino/go-rest-api-III/internal/entity"
// 	"github.com/stretchr/testify/assert"
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// func TestCreateUser(t *testing.T) {
// 	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	db.AutoMigrate(&entity.User{})

// 	user, _ := entity.NewUser("Mateus", "m@gmail.com", "123456789", "admin")
// 	userDB := NewUserDb(db)

// 	err = userDB.CreateUser(user)

// 	assert.Nil(t, err)

// 	var userFound entity.User

// 	err = db.First(&userFound, "id=?", user.ID).Error
// 	assert.Nil(t, err)
// 	assert.Equal(t, user.ID, userFound.ID)
// 	assert.Equal(t, user.Name, userFound.Name)
// 	assert.Equal(t, user.Email, userFound.Email)
// 	assert.NotNil(t, userFound.Password)
// }

// func TestFindByEmail(t *testing.T) {
// 	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	db.AutoMigrate(&entity.User{})

// 	user, _ := entity.NewUser("Mateus", "m@gmail.com", "123456789", "manager")
// 	userDB := NewUserDb(db)

// 	err = userDB.CreateUser(user)

// 	assert.Nil(t, err)

// 	// var userFound entity.User

// 	userFound, err := userDB.FindUserByEmail(user.Email)

// 	assert.Nil(t, err)
// 	assert.Equal(t, user.ID, userFound.ID)
// 	assert.Equal(t, user.Name, userFound.Name)
// 	assert.Equal(t, user.Email, userFound.Email)
// 	assert.NotNil(t, userFound.Password)
// }

// func TestFindUserById(t *testing.T) {
// 	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	db.AutoMigrate(&entity.User{})

// 	user, err := entity.NewUser("Mateus", "m@gmail.com", "123456789", "admin")

// 	assert.NoError(t, err)

// 	userDB := NewUserDb(db)

// 	err = userDB.CreateUser(user)

// 	assert.NoError(t, err)

// 	assert.NotEmpty(t, user.ID)

// 	userFound, err := userDB.FindUserById(user.ID.String())

// 	assert.NoError(t, err)

// 	assert.Equal(t, userFound.ID, user.ID)
// 	assert.Equal(t, userFound.Name, user.Name)
// 	assert.Equal(t, userFound.Email, user.Email)

// }

// func TestUpdateUser(t *testing.T) {
// 	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	db.AutoMigrate(&entity.User{})

// 	user, err := entity.NewUser("Mateus", "m@gmail.com", "123456789", "customer")

// 	assert.NoError(t, err)

// 	userDB := NewUserDb(db)

// 	err = userDB.CreateUser(user)

// 	assert.NoError(t, err)

// 	assert.NotEmpty(t, user.ID)
// 	user.Name = "Lucena"
// 	user.Role = "admin"
// 	err = userDB.UpdateUser(user)

// 	productFound, err := userDB.FindUserById(user.ID.String())
// 	assert.NoError(t, err)
// 	assert.Equal(t, "Lucena", productFound.Name)
// 	assert.Equal(t, "admin", productFound.Role)
// }
