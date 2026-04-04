package repository

import (
	"github.com/project/wayt/internal/model"
	"gorm.io/gorm"
)

type QRCodeRepository interface {
	Create(qr *model.QRCode) error
	FindByToken(token string) (*model.QRCode, error)
	FindActiveByBranch(branchID uint) ([]model.QRCode, error)
	DeactivateByBranch(branchID uint) error
}

type qrCodeRepository struct {
	db *gorm.DB
}

func NewQRCodeRepository(db *gorm.DB) QRCodeRepository {
	return &qrCodeRepository{db: db}
}

func (r *qrCodeRepository) Create(qr *model.QRCode) error {
	return r.db.Create(qr).Error
}

func (r *qrCodeRepository) FindByToken(token string) (*model.QRCode, error) {
	var qr model.QRCode
	err := r.db.Preload("Branch").Where("token = ?", token).First(&qr).Error
	if err != nil {
		return nil, err
	}
	return &qr, nil
}

func (r *qrCodeRepository) FindActiveByBranch(branchID uint) ([]model.QRCode, error) {
	var qrs []model.QRCode
	err := r.db.Where("branch_id = ? AND is_active = true AND expired_at > NOW()", branchID).
		Find(&qrs).Error
	return qrs, err
}

func (r *qrCodeRepository) DeactivateByBranch(branchID uint) error {
	return r.db.Model(&model.QRCode{}).
		Where("branch_id = ?", branchID).
		Update("is_active", false).Error
}
