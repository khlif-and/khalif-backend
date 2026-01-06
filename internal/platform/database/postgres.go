package database

import (
	"fmt"
	"os"
	"path/filepath"

	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/platform/config"
	appLogger "khalif-backend/internal/platform/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})

	if err != nil {
		appLogger.Log.Fatal("failed to connect to database", zap.Error(err))
	}

	appLogger.Log.Info("Database connected successfully")

	// Run Migrations
	runMigrations(DB)

	return DB
}

func runMigrations(db *gorm.DB) {
	appLogger.Log.Info("Running Auto-Migrations...")

	err := db.AutoMigrate(&domain.Admin{}, &domain.AdminAuditLog{}, &domain.RefreshToken{})
	if err != nil {
		appLogger.Log.Fatal("failed to auto-migrate tables", zap.Error(err))
	}
	appLogger.Log.Info("Table migration complete.")

	migrationDir := "migrations/sql"
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		appLogger.Log.Error("Could not read migration directory", zap.String("dir", migrationDir), zap.Error(err))
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			path := filepath.Join(migrationDir, file.Name())
			content, err := os.ReadFile(path)
			if err != nil {
				appLogger.Log.Error("Failed to read migration file", zap.String("file", file.Name()), zap.Error(err))
				continue
			}

			sqlQuery := string(content)
			
			if err := db.Exec(sqlQuery).Error; err != nil {
				appLogger.Log.Error("Migration File Failed", zap.String("file", file.Name()), zap.Error(err))
			} 
		}
	}
}
