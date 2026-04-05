package model

import "time"

type AdminRole string

const (
	RoleSuperAdmin AdminRole = "superadmin"
	RoleAdmin      AdminRole = "admin"
)

type AdminUser struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"      json:"id"`
	Username  string    `gorm:"size:100;not null;uniqueIndex" json:"username"`
	Role      AdminRole `gorm:"type:admin_role;default:'admin'" json:"role"`
	Password  string    `gorm:"size:255;not null"             json:"-"`
	BranchID  *uint     `gorm:"index"                         json:"branch_id,omitempty"`
	CreatedAt time.Time `                                     json:"created_at"`
	UpdatedAt time.Time `                                     json:"updated_at"`

	Branch *Branch `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}

func (AdminUser) TableName() string { return "admin_users" }
