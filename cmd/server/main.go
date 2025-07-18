package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
	"github.com/mateusfaustino/go-rest-api-III/configs"
	_ "github.com/mateusfaustino/go-rest-api-III/docs"
	"github.com/mateusfaustino/go-rest-api-III/internal/entity"
	"github.com/mateusfaustino/go-rest-api-III/internal/infra/database"
	seed "github.com/mateusfaustino/go-rest-api-III/internal/infra/database/seeds"
	"github.com/mateusfaustino/go-rest-api-III/internal/infra/webserver/handlers"
	"github.com/mateusfaustino/go-rest-api-III/internal/infra/webserver/middlewares"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// Carrega variáveis do .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Erro ao carregar .env, usando variáveis do sistema")
	}

	// Carrega configurações
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Obtém as variáveis de ambiente
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// String de conexão MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar no banco de dados: %v", err)
	}

	fmt.Println("Conectado ao MySQL com sucesso!")

	// AutoMigrate para criar as tabelas automaticamente
	db.AutoMigrate(&entity.Role{}, &entity.Product{}, &entity.User{})
	seed.SeedRoles(db)
	seed.SeedProducts(db)

	productdb := database.NewProductDB(db)
	userdb := database.NewUserDb(db)
	roledb := database.NewRoleDB(db)

	ProductHandler := handlers.NewProductHandler(productdb)
	UserHandler := handlers.NewUserHandler(userdb, roledb, cfg.TokenAuth, cfg.JwtExpiresIn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Rotas não autenticadas

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", UserHandler.GetJWT)
		r.Post("/register", UserHandler.CreateUser)
	})

	r.Route("/product", func(r chi.Router) {
		r.Get("/", ProductHandler.GetProducts)
		r.Get("/{id}", ProductHandler.GetProduct)
	})

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/docs/doc.json")))
	// Grupo para usuários autenticados
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(cfg.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Route("/user", func(r chi.Router) {
			r.Get("/profile", UserHandler.ShowOwnProfile)
			r.Put("/profile", UserHandler.UpdateOwnProfile)
			r.Get("/{id}", UserHandler.GetUserById)
		})

		r.Route("/admin", func(r chi.Router) {
			r.Use(middlewares.RoleMiddleware(roledb, "manager", "admin"))
			r.Route("/product", func(r chi.Router) {
				r.Post("/", ProductHandler.CreateProduct)
				r.Put("/{id}", ProductHandler.UpdateProduct)
				r.Delete("/{id}", ProductHandler.DeleteProduct)
			})
		})
	})

	port := os.Getenv("WEB_SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Servidor rodando na porta:", port)
	http.ListenAndServe(":"+port, r)
}
