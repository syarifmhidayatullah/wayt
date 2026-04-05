package repository

import (
	"github.com/project/wayt/internal/model"
	"gorm.io/gorm"
)

type CounterRepository interface {
	Create(counter *model.Counter) error
	FindAll() ([]model.Counter, error)
	FindByID(id uint) (*model.Counter, error)
	FindByBranch(branchID uint) ([]model.Counter, error)
	Update(counter *model.Counter) error
	Delete(id uint) error
	IncrementCurrentNumber(id uint) error
	ResetNumbers(id uint) error
}

type counterRepository struct {
	db *gorm.DB
}

func NewCounterRepository(db *gorm.DB) CounterRepository {
	return &counterRepository{db: db}
}

func (r *counterRepository) Create(counter *model.Counter) error {
	return r.db.Create(counter).Error
}

func (r *counterRepository) FindAll() ([]model.Counter, error) {
	var counters []model.Counter
	err := r.db.Where("deleted_at IS NULL").Order("branch_id ASC, id ASC").Find(&counters).Error
	return counters, err
}

func (r *counterRepository) FindByID(id uint) (*model.Counter, error) {
	var counter model.Counter
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&counter).Error
	if err != nil {
		return nil, err
	}
	return &counter, nil
}

func (r *counterRepository) FindByBranch(branchID uint) ([]model.Counter, error) {
	var counters []model.Counter
	err := r.db.Where("branch_id = ? AND deleted_at IS NULL", branchID).
		Order("id ASC").Find(&counters).Error
	return counters, err
}

func (r *counterRepository) Update(counter *model.Counter) error {
	return r.db.Save(counter).Error
}

func (r *counterRepository) Delete(id uint) error {
	return r.db.Model(&model.Counter{}).Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *counterRepository) IncrementCurrentNumber(id uint) error {
	return r.db.Model(&model.Counter{}).Where("id = ?", id).
		UpdateColumn("current_number", gorm.Expr("current_number + 1")).Error
}

func (r *counterRepository) ResetNumbers(id uint) error {
	return r.db.Model(&model.Counter{}).Where("id = ?", id).
		Updates(map[string]interface{}{"current_number": 0, "last_number": 0}).Error
}
