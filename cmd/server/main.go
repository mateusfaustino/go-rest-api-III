package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mateusfaustino/go-rest-api-III/configs"
	"github.com/mateusfaustino/go-rest-api-III/internal/entity"
	"github.com/mateusfaustino/go-rest-api-III/internal/infra/database"
	"github.com/mateusfaustino/go-rest-api-III/internal/infra/webserver/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	fmt.Println(config.DBDriver)

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&entity.Product{}, &entity.User{})
	productdb := database.NewProductDB(db)
	ProductHandler := handlers.NewProductHandler(productdb)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/product", ProductHandler.GetProducts)
	r.Post("/product", ProductHandler.CreateProduct)
	r.Get("/product/{id}", ProductHandler.GetProduct)
	r.Put("/product/{id}", ProductHandler.UpdateProduct)
	r.Delete("/product/{id}", ProductHandler.DeleteProduct)

	http.ListenAndServe(":8080", r)
}
