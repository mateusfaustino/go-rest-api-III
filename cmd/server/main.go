package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/mateusfaustino/go-rest-api-III/configs"
	"github.com/mateusfaustino/go-rest-api-III/internal/entity"
	"github.com/mateusfaustino/go-rest-api-III/internal/infra/database"
	"github.com/mateusfaustino/go-rest-api-III/internal/infra/webserver/handlers"
	"github.com/mateusfaustino/go-rest-api-III/internal/infra/webserver/middlewares"
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
	userdb := database.NewUserDb(db)

	ProductHandler := handlers.NewProductHandler(productdb)
	UserHandler := handlers.NewUserHandler(userdb, config.TokenAuth, config.JwtExpiresIn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Rotas não autenticadas
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", UserHandler.GetJWT)
		r.Post("/register", UserHandler.CreateUser)
	})

	// Grupo para usuários autenticados
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)

		// Rotas autenticadas sem restrição de role
		r.Route("/user", func(r chi.Router) {
			r.Get("/profile", UserHandler.ShowOwnProfile)
			r.Put("/profile", UserHandler.UpdateOwnProfile)
			r.Get("/{id}", UserHandler.GetUserById)
		})

		r.Route("/product", func(r chi.Router) {
			r.Get("/", ProductHandler.GetProducts)
			r.Get("/{id}", ProductHandler.GetProduct)
		})

		// Rotas protegidas (manager)
		r.Route("/manager", func(r chi.Router) {
			r.Use(middlewares.RoleMiddleware("manager"))
			r.Get("/", UserHandler.TestManager)
			// Adicione outras rotas específicas para managers aqui
		})

		// Rotas protegidas (admin)
		r.Route("/admin", func(r chi.Router) {
			r.Use(middlewares.RoleMiddleware("manager", "admin"))
			r.Get("/", UserHandler.TestAdmin)
			// Adicione outras rotas específicas para admins aqui

			r.Route("/product", func(r chi.Router) {
				r.Post("/", ProductHandler.CreateProduct)
				r.Put("/{id}", ProductHandler.UpdateProduct)
				r.Delete("/{id}", ProductHandler.DeleteProduct)
			})
		})
	})

	http.ListenAndServe(":8080", r)
}
