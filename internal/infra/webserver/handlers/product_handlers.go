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
	entityPkg "github.com/mateusfaustino/go-rest-api-III/pkg/entity"
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

	w.Header().Set("Cotent-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (ph *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() // Garante que qualquer recurso seja fechado
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, `{"error": "missing product ID"}`, http.StatusBadRequest)
		return
	}

	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	product.ID, err = entityPkg.ParseID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": %s}`, err.Error()), http.StatusInternalServerError)
		return
	}

	_, err = ph.ProductDB.FindProductByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error": "product not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	err = ph.ProductDB.UpdateProduct(&product)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": %s}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

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
