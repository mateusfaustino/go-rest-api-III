package migrations

import (
	"log"

	"github.com/mateusfaustino/go-rest-api-III/internal/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&entity.Role{}, &entity.User{})
	if err != nil {
		log.Fatalf("Erro ao rodar migração: %v", err)
	}
}
