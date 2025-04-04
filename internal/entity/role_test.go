package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRole(t *testing.T) {
	role, err := NewRole("manager")
	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.NotEmpty(t, role.ID)
	assert.Equal(t, "manager", role.Name)
}
