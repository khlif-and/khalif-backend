package admin

import (
	"errors"
	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"

	"gorm.io/gorm"
)

type AdminRepo struct {
	db *gorm.DB
}

func NewAdminRepo(db *gorm.DB) ports.AdminRepository {
	return &AdminRepo{db: db}
}

func (r *AdminRepo) Create(admin *domain.Admin) error {
	return r.db.Create(admin).Error
}

func (r *AdminRepo) FindByEmail(email string) (*domain.Admin, error) {
	var admin domain.Admin
	if err := r.db.Where("email = ?", email).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found, let service handle it
		}
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepo) FindByID(id uint) (*domain.Admin, error) {
	var admin domain.Admin
	if err := r.db.First(&admin, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepo) FindByUsername(username string) (*domain.Admin, error) {
	var admin domain.Admin
	if err := r.db.Where("username = ?", username).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepo) Update(admin *domain.Admin) error {
	return r.db.Save(admin).Error
}

func (r *AdminRepo) Count() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Admin{}).Count(&count).Error
	return count, err
}
