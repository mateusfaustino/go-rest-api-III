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
	_, err := pdb.FindProductByID(product.ID.String())

	if err != nil {
		return err
	}

	return pdb.DB.Save(product).Error
}

func (pdb *ProductDB) DeleteProduct(id string) error {
	product, err := pdb.FindProductByID(id)

	if err != nil {
		return err
	}

	return pdb.DB.Delete(product).Error
}

func (pdb *ProductDB) FindAllProducts(page, limit int, sort string) ([]entity.Product, error) {

	var products []entity.Product
	var err error
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if page != 0 && limit != 0 {
		err = pdb.DB.Limit(limit).Offset((page - 1) * limit).Order("created_at " + sort).Find(&products).Error
		return products, err
	}

	err = pdb.DB.Order("created_at " + sort).Find(&products).Error
	return products, err

}
