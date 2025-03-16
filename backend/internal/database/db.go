package database

import (
	"fmt"
	"log"
	"os"

	"github.com/thejpness/ArcadiaGo/internal/models" // ✅ Import models
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ✅ Initialize Database Connection
func InitDB() {
	dsn := os.Getenv("DATABASE_URL") // Load DSN from environment variable
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL is not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	fmt.Println("✅ Database connected successfully")
}

// ✅ Run Auto-Migrations
func Migrate() {
	if DB == nil {
		log.Fatal("❌ Database not initialized")
	}

	err := DB.AutoMigrate(
		&models.User{},            // ✅ Correctly reference models from models package
		&models.UserEmailChange{}, // ✅ Correctly reference models from models package
		&models.UserSession{},     // ✅ Correctly reference models from models package
	)

	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	log.Println("✅ Database migration completed")
}
