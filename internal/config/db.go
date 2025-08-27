package config

import (
	"interview-tracker/internal/pkg/logs"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// AutoMigrate
	migrationDB()

	// ----------------------------Database Connection---------------------------------
	logs.Logger.Printf("database| url: %s", os.Getenv("DATABASE_URL"))
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	logs.Logger.Printf("database| Connected to PostgreSQL successfully!")
	DB = db
}
