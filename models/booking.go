package models

import "gorm.io/gorm"

type Booking struct {
	gorm.Model
	BookingCode string `json:"bookingCode"`
	Phone       string `json:"phone"`

	UserID int  `json:"userId"`
	User   User `gorm:"foreignKey:user_id" json:"-"`

	EventID int   `json:"eventId"`
	Event   Event `gorm:"foreignKey:eventId" json:"event"`
}
