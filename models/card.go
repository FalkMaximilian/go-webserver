package models

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	Front string `gorm:"not null;size:255" json:"front"`
	Back  string `gorm:"not null;size:255" json:"back"`
}
