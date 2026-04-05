package repository

import (
	"github.com/project/wayt/internal/model"
	"gorm.io/gorm"
)

type QueueRepository interface {
	Create(queue *model.Queue) error
	FindByID(id uint) (*model.Queue, error)
	FindByQRToken(token string) (*model.Queue, error)
	FindActiveByCounter(counterID uint) ([]model.Queue, error)
	CountWaitingAhead(counterID uint, queueID uint) (int64, error)
	UpdateStatus(id uint, status model.QueueStatus) error
	FindNextWaiting(counterID uint) (*model.Queue, error)
	ExpireByCounter(counterID uint) error
}

type queueRepository struct {
	db *gorm.DB
}

func NewQueueRepository(db *gorm.DB) QueueRepository {
	return &queueRepository{db: db}
}

func (r *queueRepository) Create(queue *model.Queue) error {
	return r.db.Create(queue).Error
}

func (r *queueRepository) FindByID(id uint) (*model.Queue, error) {
	var queue model.Queue
	err := r.db.Preload("Branch").Preload("Counter").Where("id = ?", id).First(&queue).Error
	if err != nil {
		return nil, err
	}
	return &queue, nil
}

func (r *queueRepository) FindByQRToken(token string) (*model.Queue, error) {
	var queue model.Queue
	err := r.db.Preload("Branch").Preload("Counter").
		Where("qr_token = ? AND status IN ('waiting','called')", token).
		First(&queue).Error
	if err != nil {
		return nil, err
	}
	return &queue, nil
}

func (r *queueRepository) FindActiveByCounter(counterID uint) ([]model.Queue, error) {
	var queues []model.Queue
	err := r.db.Where("counter_id = ? AND status IN ('waiting','called')", counterID).
		Order("id ASC").Find(&queues).Error
	return queues, err
}

func (r *queueRepository) CountWaitingAhead(counterID uint, queueID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Queue{}).
		Where("counter_id = ? AND id < ? AND status = 'waiting'", counterID, queueID).
		Count(&count).Error
	return count, err
}

func (r *queueRepository) UpdateStatus(id uint, status model.QueueStatus) error {
	return r.db.Model(&model.Queue{}).Where("id = ?", id).
		Update("status", status).Error
}

func (r *queueRepository) FindNextWaiting(counterID uint) (*model.Queue, error) {
	var queue model.Queue
	err := r.db.Where("counter_id = ? AND status = 'waiting'", counterID).
		Order("id ASC").First(&queue).Error
	if err != nil {
		return nil, err
	}
	return &queue, nil
}

func (r *queueRepository) ExpireByCounter(counterID uint) error {
	return r.db.Model(&model.Queue{}).
		Where("counter_id = ? AND status IN ('waiting','called')", counterID).
		Update("status", model.QueueStatusExpired).Error
}
