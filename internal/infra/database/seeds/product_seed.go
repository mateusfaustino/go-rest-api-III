package seed

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/mateusfaustino/go-rest-api-III/internal/entity"
	"github.com/mateusfaustino/go-rest-api-III/internal/infra/database"
	"gorm.io/gorm"
)

// Lista de nomes de produtos para geração aleatória
var productNames = []string{
	"Smartphone", "Laptop", "Tablet", "Smartwatch", "Headphones",
	"Camera", "Gaming Console", "TV", "Refrigerator", "Washing Machine",
	"Microwave", "Coffee Maker", "Blender", "Vacuum Cleaner", "Air Conditioner",
	"Fan", "Toaster", "Iron", "Hair Dryer", "Electric Shaver",
}

// Lista de categorias para geração aleatória
var productCategories = []string{
	"Electronics", "Home Appliances", "Computing", "Gaming", "Audio",
	"Kitchen", "Cleaning", "Personal Care", "Office", "Sports",
}

func SeedProducts(db *gorm.DB) {
	productDB := database.NewProductDB(db)

	// Verifica se já existem produtos no banco
	var count int64
	if err := db.Model(&entity.Product{}).Count(&count).Error; err != nil {
		log.Printf("Erro ao verificar produtos existentes: %v\n", err)
		return
	}

	if count > 0 {
		fmt.Println("Já existem produtos no banco de dados. Pulando seed...")
		return
	}

	// Inicializa o gerador de números aleatórios
	rand.Seed(time.Now().UnixNano())

	// Cria 20 produtos aleatórios
	for i := 0; i < 20; i++ {
		// Gera um nome aleatório
		nameIndex := rand.Intn(len(productNames))
		categoryIndex := rand.Intn(len(productCategories))
		name := fmt.Sprintf("%s %s", productNames[nameIndex], productCategories[categoryIndex])

		// Gera um preço aleatório entre 100 e 10000
		price := rand.Float64()*9900 + 100

		// Cria o produto
		product, err := entity.NewProduct(name, price)
		if err != nil {
			log.Printf("Erro ao criar produto '%s': %v\n", name, err)
			continue
		}

		// Salva no banco
		if err := productDB.CreateProduct(product); err != nil {
			log.Printf("Erro ao salvar produto '%s' no banco: %v\n", name, err)
		} else {
			fmt.Printf("Produto '%s' criado com sucesso!\n", name)
		}
	}
}
