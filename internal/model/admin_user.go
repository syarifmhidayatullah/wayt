package model

import "time"

type AdminRole string

const (
	RoleSuperAdmin AdminRole = "superadmin"
	RoleAdmin      AdminRole = "admin"
)

type AdminUser struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"size:100;not null;uniqueIndex" json:"username"`
	Role      AdminRole `gorm:"type:enum('superadmin','admin');default:'admin'" json:"role"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (AdminUser) TableName() string { return "admin_users" }
