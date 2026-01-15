package database

import (
	"fmt"
	"path/filepath"

	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/platform/config"
	appLogger "khalif-backend/internal/platform/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	// Run GORM AutoMigrate for model sync
	runAutoMigrate(DB)

	// Run SQL Migrations with golang-migrate (tracked)
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)
	runMigrations(databaseURL)

	return DB
}

func runAutoMigrate(db *gorm.DB) {
	appLogger.Log.Info("Running GORM AutoMigrate...")

	// List of models to migrate - GORM handles struct-based migrations
	models := []interface{}{
		&domain.Admin{},
		&domain.AdminAuditLog{},
		&domain.RefreshToken{},
		&domain.User{},
		&domain.UserAuditLog{},
		&domain.UserRefreshToken{},
		&domain.OTPToken{},
		&domain.PasswordResetToken{},
		&domain.MoodCategory{},
		&domain.Ustadz{},
		&domain.Audio{},
		&domain.Like{},
		&domain.ListeningHistory{},
		&domain.Playlist{},
		&domain.PlaylistAudio{},
		&domain.PlaylistLike{},
		&domain.Hadist{},
		&domain.HadistLike{},
		&domain.HadistBookmark{},
		&domain.Doa{},
		&domain.DoaLike{},
		&domain.DoaBookmark{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		appLogger.Log.Error("AutoMigrate failed", zap.Error(err))
	} else {
		appLogger.Log.Info("AutoMigrate completed successfully")
	}
}

func runMigrations(databaseURL string) {
	appLogger.Log.Info("Running SQL Migrations with golang-migrate...")

	// Get absolute path to migrations directory
	migrationPath, err := filepath.Abs("migrations/sql")
	if err != nil {
		appLogger.Log.Error("Failed to get migration path", zap.Error(err))
		return
	}

	sourceURL := fmt.Sprintf("file://%s", migrationPath)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		appLogger.Log.Error("Failed to create migrate instance", zap.Error(err))
		return
	}
	defer m.Close()

	// Run all pending migrations
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			appLogger.Log.Info("No new migrations to apply")
		} else {
			appLogger.Log.Error("Migration failed", zap.Error(err))
		}
		return
	}

	version, dirty, _ := m.Version()
	appLogger.Log.Info("Migrations applied successfully",
		zap.Uint("version", version),
		zap.Bool("dirty", dirty))
}
