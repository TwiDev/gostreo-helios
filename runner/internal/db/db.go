package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := os.Getenv("DATABASE_URL") // Exemple: "host=localhost user=postgres password=xxx dbname=workflow port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto-migration des modèles
	err = db.AutoMigrate(
		&Workflow{},
		&WorkflowNode{},
		&NodeParameter{},
		&ExecutionLog{},
		&RunHistory{},
		&Event{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database schema: %v", err)
	}

	DB = db
	fmt.Println("Database connected and migrated")
}
