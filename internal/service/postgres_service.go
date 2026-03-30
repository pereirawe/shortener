package service

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/pereirawe/shortener/internal/config"
	"github.com/pereirawe/shortener/internal/model"
)

// PostgresService handles all PostgreSQL operations via GORM
type PostgresService struct {
	DB *gorm.DB
}

// NewPostgresService creates a new PostgresService and runs auto migrations
func NewPostgresService(cfg *config.Config) (*PostgresService, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=America/Sao_Paulo",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.PostgresPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Run auto migrations
	if err := db.AutoMigrate(&model.URL{}); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Printf("Connected to PostgreSQL at %s:%d (db=%s)", cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB)
	return &PostgresService{DB: db}, nil
}

// FindByShortCode fetches a URL record from the database by short code
func (ps *PostgresService) FindByShortCode(shortCode string) (*model.URL, error) {
	var url model.URL
	result := ps.DB.Where("short_code = ?", shortCode).First(&url)
	if result.Error != nil {
		return nil, result.Error
	}
	return &url, nil
}

// Create inserts a new URL record into the database
func (ps *PostgresService) Create(url *model.URL) error {
	return ps.DB.Create(url).Error
}

// IncrementClicks atomically increments the click counter for a given short code
func (ps *PostgresService) IncrementClicks(shortCode string) {
	ps.DB.Model(&model.URL{}).
		Where("short_code = ?", shortCode).
		UpdateColumn("clicks", gorm.Expr("clicks + ?", 1))
}

// ExistsByShortCode checks if a short code is already taken
func (ps *PostgresService) ExistsByShortCode(shortCode string) (bool, error) {
	var count int64
	result := ps.DB.Model(&model.URL{}).Where("short_code = ?", shortCode).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
