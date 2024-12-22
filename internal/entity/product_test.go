package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	product, err := NewProduct("pruduct", 9.99)
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.NotEmpty(t, product.ID)
	assert.Equal(t, product.Name, "pruduct")
	assert.Equal(t, product.Price, 9.99)

}

func TestProductWhenNameIsRequired(t *testing.T) {
	product, err := NewProduct("", 9.99)
	assert.Nil(t, product)
	assert.Equal(t, err, ErrNameIsRequired)
}
func TestProductWhenPriceIsRequired(t *testing.T) {
	product, err := NewProduct("product", 0.00)
	assert.Nil(t, product)
	assert.Equal(t, err, ErrPriceIsRequired)
}
func TestProductWhenInvalidPrice(t *testing.T) {
	product, err := NewProduct("product", -1.1)
	assert.Nil(t, product)
	assert.Equal(t, err, ErrInvalidPrice)
}
func TestProductValidate(t *testing.T) {
	product, err := NewProduct("product", 1.1)
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.Nil(t, product.ValidateProduct())
}
