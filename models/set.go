package models

import (
	"time"
)

type Set struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"-"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Name      string    `gorm:"not null;size:255" json:"name"`
	Cards     []Card    `gorm:"many2many:set_cards" json:"-"`
}
