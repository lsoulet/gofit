package db

import (
	"log"

	"github.com/lsoulet/gofit/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() error {
	dsn := "host=localhost user=postgres password=postgres dbname=gofitdb port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erreur lors de la connexion à la base de données : %v", err)
		return err
	}

	err = DB.AutoMigrate(&models.User{}, &models.DailyMenu{}, &models.Meal{})
	if err != nil {
		return err
	}

	return nil
}
