package model

import "time"

type Branch struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string     `gorm:"size:100;not null"        json:"name"`
	Prefix        string     `gorm:"size:10;not null"         json:"prefix"`
	IsActive      bool       `gorm:"default:true"             json:"is_active"`
	CurrentNumber int        `gorm:"default:0"                json:"current_number"`
	LastNumber    int        `gorm:"default:0"                json:"last_number"`
	CreatedAt     time.Time  `                                json:"created_at"`
	UpdatedAt     time.Time  `                                json:"updated_at"`
	DeletedAt     *time.Time `gorm:"index"                    json:"deleted_at,omitempty"`
}

func (Branch) TableName() string { return "branches" }
