package repository

import (
	"github.com/project/wayt/internal/model"
	"gorm.io/gorm"
)

type AdminUserRepository interface {
	FindByUsername(username string) (*model.AdminUser, error)
	FindByID(id uint) (*model.AdminUser, error)
	FindAll() ([]model.AdminUser, error)
	Create(user *model.AdminUser) error
	Update(user *model.AdminUser) error
	Delete(id uint) error
	ExistsAny() (bool, error)
}

type adminUserRepository struct {
	db *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) AdminUserRepository {
	return &adminUserRepository{db: db}
}

func (r *adminUserRepository) FindByUsername(username string) (*model.AdminUser, error) {
	var user model.AdminUser
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *adminUserRepository) FindByID(id uint) (*model.AdminUser, error) {
	var user model.AdminUser
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *adminUserRepository) FindAll() ([]model.AdminUser, error) {
	var users []model.AdminUser
	err := r.db.Order("created_at ASC").Find(&users).Error
	return users, err
}

func (r *adminUserRepository) Create(user *model.AdminUser) error {
	return r.db.Create(user).Error
}

func (r *adminUserRepository) Update(user *model.AdminUser) error {
	return r.db.Save(user).Error
}

func (r *adminUserRepository) Delete(id uint) error {
	return r.db.Delete(&model.AdminUser{}, id).Error
}

func (r *adminUserRepository) ExistsAny() (bool, error) {
	var count int64
	err := r.db.Model(&model.AdminUser{}).Count(&count).Error
	return count > 0, err
}
