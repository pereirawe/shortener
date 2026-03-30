package model

import (
	"time"

	"gorm.io/gorm"
)

// URL represents a shortened URL record in the database
type URL struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	ShortCode   string         `gorm:"uniqueIndex;not null;size:20" json:"short_code"`
	OriginalURL string         `gorm:"not null;type:text" json:"original_url"`
	Clicks      int64          `gorm:"default:0" json:"clicks"`
}
