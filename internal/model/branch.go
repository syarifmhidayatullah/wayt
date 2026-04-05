package model

import "time"

type Branch struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"size:100;not null"        json:"name"`
	IsActive  bool       `gorm:"default:true"             json:"is_active"`
	CreatedAt time.Time  `                                json:"created_at"`
	UpdatedAt time.Time  `                                json:"updated_at"`
	DeletedAt *time.Time `gorm:"index"                    json:"deleted_at,omitempty"`
}

func (Branch) TableName() string { return "branches" }
