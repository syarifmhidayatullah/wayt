package model

import "time"

type QueueStatus string

const (
	QueueStatusWaiting QueueStatus = "waiting"
	QueueStatusCalled  QueueStatus = "called"
	QueueStatusDone    QueueStatus = "done"
	QueueStatusExpired QueueStatus = "expired"
)

type Queue struct {
	ID          uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchID    uint        `gorm:"not null;index"           json:"branch_id"`
	QRToken     string      `gorm:"size:36;not null;index"   json:"qr_token"`
	QueueNumber string      `gorm:"size:20;not null"         json:"queue_number"`
	Status      QueueStatus `gorm:"type:enum('waiting','called','done','expired');default:'waiting'" json:"status"`
	CreatedAt   time.Time   `                                json:"created_at"`
	UpdatedAt   time.Time   `                                json:"updated_at"`

	Branch *Branch `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}

func (Queue) TableName() string { return "queues" }
