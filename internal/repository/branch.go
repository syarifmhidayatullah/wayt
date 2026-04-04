package repository

import (
	"github.com/project/wayt/internal/model"
	"gorm.io/gorm"
)

type BranchRepository interface {
	Create(branch *model.Branch) error
	FindAll() ([]model.Branch, error)
	FindByID(id uint) (*model.Branch, error)
	Update(branch *model.Branch) error
	Delete(id uint) error
	IncrementCurrentNumber(id uint) error
	ResetNumbers(id uint) error
}

type branchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) BranchRepository {
	return &branchRepository{db: db}
}

func (r *branchRepository) Create(branch *model.Branch) error {
	return r.db.Create(branch).Error
}

func (r *branchRepository) FindAll() ([]model.Branch, error) {
	var branches []model.Branch
	err := r.db.Where("deleted_at IS NULL").Find(&branches).Error
	return branches, err
}

func (r *branchRepository) FindByID(id uint) (*model.Branch, error) {
	var branch model.Branch
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&branch).Error
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r *branchRepository) Update(branch *model.Branch) error {
	return r.db.Save(branch).Error
}

func (r *branchRepository) Delete(id uint) error {
	return r.db.Model(&model.Branch{}).Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *branchRepository) IncrementCurrentNumber(id uint) error {
	return r.db.Model(&model.Branch{}).Where("id = ?", id).
		UpdateColumn("current_number", gorm.Expr("current_number + 1")).Error
}

func (r *branchRepository) ResetNumbers(id uint) error {
	return r.db.Model(&model.Branch{}).Where("id = ?", id).
		Updates(map[string]interface{}{"current_number": 0, "last_number": 0}).Error
}
