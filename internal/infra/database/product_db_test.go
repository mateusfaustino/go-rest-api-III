package database

import (
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

	product, _ := entity.NewProduct("Blusa", 9.99)
	productDB := NewProductDB(db)

	err = productDB.CreateProduct(product)

	assert.Nil(t, err)
	
	var productFound entity.Product
	
	err = db.First(&productFound, "id=?", product.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, product.Name, productFound.Name)
	assert.Equal(t, product.Price, productFound.Price)

}
