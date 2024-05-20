package models

import "gorm.io/gorm"

type Set struct {
	gorm.Model
	UserID uint
	Cards  []Card `gorm:"many2many:set_cards"`
}
