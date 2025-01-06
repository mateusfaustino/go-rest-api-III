package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/mateusfaustino/go-rest-api-III/internal/entity"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Blusa", 9.99)

	assert.NoError(t, err)

	productDB := NewProductDB(db)

	err = productDB.CreateProduct(product)

	assert.NoError(t, err)

	var productFound entity.Product

	assert.NotEmpty(t, product.ID)

	err = db.First(&productFound, "id=?", product.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, product.Name, productFound.Name)
	assert.Equal(t, product.Price, productFound.Price)

}

func TestFindAllProducts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})

	for i := 1; i < 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Float64()*100)
		assert.NoError(t, err)
		db.Create(product)
	}

	productDB := NewProductDB(db)

	products, err := productDB.FindAllProducts(1, 10, "asc")

	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	products, err = productDB.FindAllProducts(2, 10, "asc")

	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 11", products[0].Name)
	assert.Equal(t, "Product 20", products[9].Name)
}

func TestFindProductById(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Blusa", 9.99)

	assert.NoError(t, err)

	productDB := NewProductDB(db)

	err = productDB.CreateProduct(product)

	assert.NoError(t, err)

	assert.NotEmpty(t, product.ID)

	productFound, err := productDB.FindProductByID(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, product.Name, productFound.Name)
	assert.Equal(t, product.ID, productFound.ID)
	assert.Equal(t, product.Price, productFound.Price)
}

func TestUpdateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Blusa", 9.99)

	assert.NoError(t, err)

	productDB := NewProductDB(db)

	err = productDB.CreateProduct(product)

	assert.NoError(t, err)

	assert.NotEmpty(t, product.ID)
	product.Name = "Blusa 2"
	err = productDB.UpdateProduct(product)

	productFound, err := productDB.FindProductByID(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, "Blusa 2", productFound.Name)
}

func TestDeleteProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Blusa", 9.99)

	assert.NoError(t, err)

	productDB := NewProductDB(db)

	err = productDB.CreateProduct(product)

	assert.NoError(t, err)

	var productFound entity.Product

	assert.NotEmpty(t, product.ID)

	err = db.First(&productFound, "id=?", product.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, product.Name, productFound.Name)
	assert.Equal(t, product.Price, productFound.Price)

	err = productDB.DeleteProduct(product.ID.String())
	assert.NoError(t, err)

	_, err = productDB.FindProductByID(product.ID.String())

	assert.Error(t, err)

}
