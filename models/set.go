package models

import "gorm.io/gorm"

type Set struct {
	gorm.Model
	UserID uint
	Name   string `gorm:"not null;size:255"`
	Cards  []Card `gorm:"many2many:set_cards"`
}
