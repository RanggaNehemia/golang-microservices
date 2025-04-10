package models

import "time"

type Price struct {
	ID        uint      `gorm:"primaryKey"`
	Value     float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
