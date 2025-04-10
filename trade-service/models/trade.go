package models

import "gorm.io/gorm"

type Trade struct {
	gorm.Model
	UserID   uint
	Price    float64
	Quantity int
}
