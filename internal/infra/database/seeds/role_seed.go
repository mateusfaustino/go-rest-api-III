package seed

import (
	"fmt"
	"log"

	"github.com/mateusfaustino/go-rest-api-III/internal/entity"
	"github.com/mateusfaustino/go-rest-api-III/internal/infra/database"
	"gorm.io/gorm"
)

// func SeedRoles(db *gorm.DB) {
// 	roleDB := database.NewRoleDB(db)

// 	roles := []string{"admin", "manager", "customer"}

// 	for _, roleName := range roles {
// 		role, err := entity.NewRole(roleName)

// 		if err != nil {
// 			log.Printf("Erro ao criar role '%s': %v\n", roleName, err)
// 		} else {
// 			fmt.Printf("Role '%s' criada com sucesso!\n", roleName)
// 		}

// 		err = roleDB.CreateRole(role)
// 		if err != nil {
// 			log.Printf("Erro ao criar role '%s': %v\n", roleName, err)
// 		} else {
// 			fmt.Printf("Role '%s' criada com sucesso!\n", roleName)
// 		}
// 	}
// }


func SeedRoles(db *gorm.DB) {
	roleDB := database.NewRoleDB(db)

	roles := []string{"admin", "manager", "customer"}

	for _, roleName := range roles {
		// Verifica se a role já existe
		exists, err := roleDB.RoleExists(roleName)
		if err != nil {
			log.Printf("Erro ao verificar a existência da role '%s': %v\n", roleName, err)
			continue
		}

		if exists {
			fmt.Printf("Role '%s' já existe. Pulando...\n", roleName)
			continue
		}

		// Criando a role
		role, err := entity.NewRole(roleName)
		if err != nil {
			log.Printf("Erro ao criar role '%s': %v\n", roleName, err)
			continue
		}

		// Salvando no banco
		if err := roleDB.CreateRole(role); err != nil {
			log.Printf("Erro ao salvar role '%s' no banco: %v\n", roleName, err)
		} else {
			fmt.Printf("Role '%s' criada com sucesso!\n", roleName)
		}
	}
}
