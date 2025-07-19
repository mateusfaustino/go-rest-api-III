package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	roleID := NewID()
	user, err := NewUser("Mateus", "m.m@gmail.com", "123456", roleID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, "Mateus", user.Name)
	assert.Equal(t, "m.m@gmail.com", user.Email)
	assert.Equal(t, roleID, user.RoleID)
	assert.NotEqual(t, "123456", user.Password)
}

func TestUser_ValidatePassword(t *testing.T) {
	roleID := NewID()
	user, err := NewUser("Mateus", "m.m@gmail.com", "123456", roleID)
	assert.NoError(t, err)
	assert.True(t, user.ValidatePassword("123456"))
	assert.False(t, user.ValidatePassword("1234567"))
	assert.NotEqual(t, "123456", user.Password)
}
