package schema

import (
	"time"
)

type User struct {
	ID        string `gorm:"type:varchar(255);primary_key;uniqueIndex;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserName  string `gorm:"not null;uniqueIndex"`
	Password  string `gorm:"not null"`
	Gender    string
	Age       int
	Email     string
	Phone     string
	Address   string
}
