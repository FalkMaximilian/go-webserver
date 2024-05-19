package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null;size:50;" validate:"required;min=3;max=50" json:"username"`
	Email    string `gorm:"uniqueIndex;not null;size:255;" validate:"required,email" json:"email"`
	Password string `gorm:"not null;" validate:"required" json:"password"`
	Sets     []Set
}

type Set struct {
	gorm.Model
	UserID uint
	Cards  []Card `gorm:"many2many:set_cards"`
}

type Card struct {
	gorm.Model
}
