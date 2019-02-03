package models

import "github.com/jinzhu/gorm"

type Ticket struct {
	gorm.Model
	IsBooked bool `gorm:"not null" json:"is_booked"`
}