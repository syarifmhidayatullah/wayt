package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"      json:"id"`
	Name      string    `gorm:"size:100;not null"             json:"name"`
	Phone     string    `gorm:"size:20;not null;uniqueIndex"  json:"phone"`
	Password  string    `gorm:"size:255;not null"             json:"-"`
	CreatedAt time.Time `                                     json:"created_at"`
	UpdatedAt time.Time `                                     json:"updated_at"`
}

func (User) TableName() string { return "users" }
