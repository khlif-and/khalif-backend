package ustadz

import (
	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"

	"gorm.io/gorm"
)

type ustadzRepo struct {
	db *gorm.DB
}

func NewUstadzRepo(db *gorm.DB) ports.UstadzRepository {
	return &ustadzRepo{db: db}
}

func (r *ustadzRepo) Create(ustadz *domain.Ustadz) error {
	return r.db.Create(ustadz).Error
}

func (r *ustadzRepo) FindByID(id uint) (*domain.Ustadz, error) {
	var ustadz domain.Ustadz
	err := r.db.First(&ustadz, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &ustadz, err
}

func (r *ustadzRepo) FindByUUID(uuid string) (*domain.Ustadz, error) {
	var ustadz domain.Ustadz
	err := r.db.Where("uuid = ?", uuid).First(&ustadz).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &ustadz, err
}

func (r *ustadzRepo) FindByName(name string) (*domain.Ustadz, error) {
	var ustadz domain.Ustadz
	err := r.db.Where("name = ?", name).First(&ustadz).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &ustadz, err
}

func (r *ustadzRepo) FindAll(page, limit int) ([]domain.Ustadz, int64, error) {
	var ustadzList []domain.Ustadz
	var total int64

	r.db.Model(&domain.Ustadz{}).Count(&total)

	offset := (page - 1) * limit
	err := r.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&ustadzList).Error

	return ustadzList, total, err
}

func (r *ustadzRepo) Update(ustadz *domain.Ustadz) error {
	return r.db.Save(ustadz).Error
}

func (r *ustadzRepo) Delete(id uint) error {
	return r.db.Delete(&domain.Ustadz{}, id).Error
}
