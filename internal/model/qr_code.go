package model

import "time"

type QRCode struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"     json:"id"`
	BranchID  uint      `gorm:"not null;index"               json:"branch_id"`
	CounterID uint      `gorm:"not null;index"               json:"counter_id"`
	Token     string    `gorm:"size:36;not null;uniqueIndex" json:"token"`
	IsActive  bool      `gorm:"default:true"                 json:"is_active"`
	ExpiredAt time.Time `gorm:"not null"                     json:"expired_at"`
	CreatedAt time.Time `                                    json:"created_at"`
	UpdatedAt time.Time `                                    json:"updated_at"`

	Branch  *Branch  `gorm:"foreignKey:BranchID"  json:"branch,omitempty"`
	Counter *Counter `gorm:"foreignKey:CounterID" json:"counter,omitempty"`
}

func (QRCode) TableName() string { return "qr_codes" }
