package hadist

import (
	"errors"

	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"

	"gorm.io/gorm"
)

type HadistRepo struct {
	db *gorm.DB
}

func NewHadistRepo(db *gorm.DB) ports.HadistRepository {
	return &HadistRepo{db: db}
}

func (r *HadistRepo) Create(hadist *domain.Hadist) error {
	return r.db.Create(hadist).Error
}

func (r *HadistRepo) FindByID(id uint) (*domain.Hadist, error) {
	var hadist domain.Hadist
	if err := r.db.First(&hadist, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &hadist, nil
}

func (r *HadistRepo) FindByUUID(uuid string) (*domain.Hadist, error) {
	var hadist domain.Hadist
	if err := r.db.Where("uuid = ?", uuid).First(&hadist).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &hadist, nil
}

func (r *HadistRepo) FindAll(page, limit int) ([]domain.Hadist, int64, error) {
	var hadists []domain.Hadist
	var total int64

	if err := r.db.Model(&domain.Hadist{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := r.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&hadists).Error; err != nil {
		return nil, 0, err
	}

	return hadists, total, nil
}

func (r *HadistRepo) FindByCategory(category string, page, limit int) ([]domain.Hadist, int64, error) {
	var hadists []domain.Hadist
	var total int64

	query := r.db.Model(&domain.Hadist{}).Where("category_hadist = ?", category)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := query.Order("nomor_hadist ASC").Offset(offset).Limit(limit).Find(&hadists).Error; err != nil {
		return nil, 0, err
	}

	return hadists, total, nil
}

func (r *HadistRepo) FindByKitab(kitab string, page, limit int) ([]domain.Hadist, int64, error) {
	var hadists []domain.Hadist
	var total int64

	query := r.db.Model(&domain.Hadist{}).Where("kitab_hadist = ?", kitab)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := query.Order("nomor_hadist ASC").Offset(offset).Limit(limit).Find(&hadists).Error; err != nil {
		return nil, 0, err
	}

	return hadists, total, nil
}

func (r *HadistRepo) FindRandom() (*domain.Hadist, error) {
	var hadist domain.Hadist
	if err := r.db.Order("RANDOM()").First(&hadist).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &hadist, nil
}

func (r *HadistRepo) Update(hadist *domain.Hadist) error {
	return r.db.Save(hadist).Error
}

func (r *HadistRepo) Delete(id uint) error {
	return r.db.Delete(&domain.Hadist{}, id).Error
}

func (r *HadistRepo) IncrementListeningCount(id uint) error {
	return r.db.Model(&domain.Hadist{}).Where("id = ?", id).
		UpdateColumn("listening_count", gorm.Expr("listening_count + ?", 1)).Error
}

// Like operations
func (r *HadistRepo) CreateLike(like *domain.HadistLike) error {
	return r.db.Create(like).Error
}

func (r *HadistRepo) FindLikeByUserAndHadist(userID, hadistID uint) (*domain.HadistLike, error) {
	var like domain.HadistLike
	if err := r.db.Where("user_id = ? AND hadist_id = ?", userID, hadistID).First(&like).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &like, nil
}

func (r *HadistRepo) DeleteLike(id uint) error {
	return r.db.Delete(&domain.HadistLike{}, id).Error
}

func (r *HadistRepo) IncrementLikeCount(hadistID uint) error {
	return r.db.Model(&domain.Hadist{}).Where("id = ?", hadistID).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
}

func (r *HadistRepo) DecrementLikeCount(hadistID uint) error {
	return r.db.Model(&domain.Hadist{}).Where("id = ?", hadistID).
		UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error
}

// Bookmark operations
func (r *HadistRepo) CreateBookmark(bookmark *domain.HadistBookmark) error {
	return r.db.Create(bookmark).Error
}

func (r *HadistRepo) FindBookmarkByUserAndHadist(userID, hadistID uint) (*domain.HadistBookmark, error) {
	var bookmark domain.HadistBookmark
	if err := r.db.Where("user_id = ? AND hadist_id = ?", userID, hadistID).First(&bookmark).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &bookmark, nil
}

func (r *HadistRepo) DeleteBookmark(id uint) error {
	return r.db.Delete(&domain.HadistBookmark{}, id).Error
}

func (r *HadistRepo) IncrementBookmarkCount(hadistID uint) error {
	return r.db.Model(&domain.Hadist{}).Where("id = ?", hadistID).
		UpdateColumn("bookmark_count", gorm.Expr("bookmark_count + ?", 1)).Error
}

func (r *HadistRepo) DecrementBookmarkCount(hadistID uint) error {
	return r.db.Model(&domain.Hadist{}).Where("id = ?", hadistID).
		UpdateColumn("bookmark_count", gorm.Expr("GREATEST(bookmark_count - 1, 0)")).Error
}

func (r *HadistRepo) GetUserBookmarks(userID uint, page, limit int) ([]domain.HadistBookmark, int64, error) {
	var bookmarks []domain.HadistBookmark
	var total int64

	if err := r.db.Model(&domain.HadistBookmark{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&bookmarks).Error; err != nil {
		return nil, 0, err
	}

	return bookmarks, total, nil
}
