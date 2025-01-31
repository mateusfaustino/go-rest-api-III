package database

import (
	"github.com/mateusfaustino/go-rest-api-III/internal/entity"
	"gorm.io/gorm"
)

type ProductDB struct {
	DB *gorm.DB
}

func NewProductDB(db *gorm.DB) *ProductDB {
	return &ProductDB{
		DB: db,
	}
}

func (pdb *ProductDB) CreateProduct(product *entity.Product) error {
	return pdb.DB.Create(product).Error
}

func (pdb *ProductDB) FindProductByID(id string) (*entity.Product, error) {
	var product entity.Product
	err := pdb.DB.First(&product, "id=?", id).Error
	return &product, err
}

func (pdb *ProductDB) UpdateProduct(product *entity.Product) error {
	return pdb.DB.Save(product).Error
}

func (pdb *ProductDB) DeleteProduct(id string) error {
	result := pdb.DB.Delete(&entity.Product{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (pdb *ProductDB) FindAllProducts(page, limit int, sort string) ([]entity.Product, error) {
	var products []entity.Product
	offset := (page - 1) * limit // Calculando o offset corretamente

	// Executando a query com paginação e ordenação seguras
	err := pdb.DB.
		Order("created_at " + sort).
		Limit(limit).
		Offset(offset).
		Find(&products).
		Error

	return products, err
}
