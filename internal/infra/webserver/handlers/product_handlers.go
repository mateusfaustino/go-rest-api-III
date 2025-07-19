package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/mateusfaustino/go-rest-api-III/internal/dto"
	"github.com/mateusfaustino/go-rest-api-III/internal/entity"
	"github.com/mateusfaustino/go-rest-api-III/internal/infra/database"
	"gorm.io/gorm"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

func (ph *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var productInput dto.CreateProductInput

	// Fechar o corpo da requisição após uso
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&productInput)

	if err != nil {
		http.Error(w, `{"error": "invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	p, err := entity.NewProduct(productInput.Name, productInput.Price)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": %s}`, err.Error()), http.StatusBadRequest)
		return
	}

	err = ph.ProductDB.CreateProduct(p)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": %s}`, err.Error()), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "product created successfully"})
}

// Create Product godoc
// @Summary Create a new product
// @Description Create a new product with the given name and price
// @Tags products
// @Accept json
// @Produce json
// @Param product body dto.CreateProductInput true "Product data"
// @Success 200 {object} map[string]string
// @Router /product [post]
// @Security ApiKeyAuth
// GetProduct godoc
// @Summary Get a product
// @Description Retrieve a product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} entity.Product
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /product/{id} [get]
func (ph *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() // Garante que qualquer recurso seja fechado
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, `{"error": "missing product ID"}`, http.StatusBadRequest)
		return
	}

	product, err := ph.ProductDB.FindProductByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error": "product not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Update product identified by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body dto.UpdateProductInput true "Product data"
// @Success 200 {object} entity.Product
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /admin/product/{id} [put]
// @Security ApiKeyAuth
func (ph *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() // Garante que qualquer recurso seja fechado
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, `{"error": "missing product ID"}`, http.StatusBadRequest)
		return
	}

	var input dto.UpdateProductInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	product, err := ph.ProductDB.FindProductByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error": "product not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	// Verifica se o campo foi setado no corpo da requisição
	if input.Name != "" {
		product.Name = input.Name
	}

	if input.Price != 0.0 {
		product.Price = input.Price
	}

	err = ph.ProductDB.UpdateProduct(product)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": %s}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /admin/product/{id} [delete]
// @Security ApiKeyAuth
func (ph *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, `{"error": "missing product ID"}`, http.StatusBadRequest)
		return
	}

	err := ph.ProductDB.DeleteProduct(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error": "product not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "product deleted successfully"})
}

// GetProducts godoc
// @Summary List products
// @Description Get products with pagination
// @Tags products
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort order"
// @Success 200 {array} entity.Product
// @Failure 500 {object} Error
// @Router /product [get]
func (ph *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")

	pageInt, err := strconv.Atoi(page)

	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 10 // Define um limite padrão
	}

	// Segurança na ordenação
	validSortOptions := map[string]bool{"asc": true, "desc": true}
	if !validSortOptions[sort] {
		sort = "asc"
	}

	products, err := ph.ProductDB.FindAllProducts(pageInt, limitInt, sort)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": %s}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
