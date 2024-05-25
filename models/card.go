package models

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	UserID uint
	Front  string `gorm:"not null;size:255" json:"front"`
	Back   string `gorm:"not null;size:255" json:"back"`
	Sets   []Set  `gorm:"many2many:set_cards"`
}
