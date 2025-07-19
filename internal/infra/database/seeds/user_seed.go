package seed

import (
	"fmt"
	"log"

	"github.com/mateusfaustino/go-rest-api-III/internal/entity"
	"github.com/mateusfaustino/go-rest-api-III/internal/infra/database"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	userDB := database.NewUserDb(db)
	roleDB := database.NewRoleDB(db)

	// Check if users already exist
	var count int64
	if err := db.Model(&entity.User{}).Count(&count).Error; err != nil {
		log.Printf("Erro ao verificar usuários existentes: %v\n", err)
		return
	}

	if count > 0 {
		fmt.Println("Já existem usuários no banco de dados. Pulando seed...")
		return
	}

	// Ensure roles exist
	roleCustomer, _ := roleDB.FindRoleByName("customer")
	roleManager, _ := roleDB.FindRoleByName("manager")
	roleAdmin, _ := roleDB.FindRoleByName("admin")

	users := []struct {
		name  string
		email string
		role  *entity.Role
	}{
		{"Customer User", "customer@example.com", roleCustomer},
		{"Manager User", "manager@example.com", roleManager},
		{"Admin User", "admin@example.com", roleAdmin},
	}

	for _, u := range users {
		if u.role == nil {
			log.Printf("Role não encontrada para o usuário %s\n", u.email)
			continue
		}
		user, err := entity.NewUser(u.name, u.email, "1234", u.role.ID)
		if err != nil {
			log.Printf("Erro ao criar usuário '%s': %v\n", u.email, err)
			continue
		}
		if err := userDB.CreateUser(user); err != nil {
			log.Printf("Erro ao salvar usuário '%s': %v\n", u.email, err)
		} else {
			fmt.Printf("Usuário '%s' criado com sucesso!\n", u.email)
		}
	}
}
