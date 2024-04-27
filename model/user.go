package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null;size:50;" validate:"required;min=3;max=50"`
	Email    string `gorm:"uniqueIndex;not null;size:255;" validate:"required,email"`
	Password string `gorm:"not null;" validate:"required"`
}
