package doa

import (
	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"

	"gorm.io/gorm"
)

type DoaRepo struct {
	db *gorm.DB
}

func NewDoaRepo(db *gorm.DB) ports.DoaRepository {
	return &DoaRepo{db: db}
}

func (r *DoaRepo) Create(doa *domain.Doa) error {
	return r.db.Create(doa).Error
}

func (r *DoaRepo) FindByID(id uint) (*domain.Doa, error) {
	var doa domain.Doa
	if err := r.db.Preload("Hadist").First(&doa, id).Error; err != nil {
		return nil, err
	}
	return &doa, nil
}

func (r *DoaRepo) FindByUUID(uuid string) (*domain.Doa, error) {
	var doa domain.Doa
	if err := r.db.Preload("Hadist").Where("uuid = ?", uuid).First(&doa).Error; err != nil {
		return nil, err
	}
	return &doa, nil
}

func (r *DoaRepo) FindAll(page, limit int) ([]domain.Doa, int64, error) {
	var doas []domain.Doa
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&domain.Doa{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Hadist").Limit(limit).Offset(offset).Order("created_at desc").Find(&doas).Error; err != nil {
		return nil, 0, err
	}

	return doas, total, nil
}

func (r *DoaRepo) FindByCategory(category string, page, limit int) ([]domain.Doa, int64, error) {
	var doas []domain.Doa
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&domain.Doa{}).Where("category_doa = ?", category)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Hadist").Limit(limit).Offset(offset).Order("created_at desc").Find(&doas).Error; err != nil {
		return nil, 0, err
	}

	return doas, total, nil
}

func (r *DoaRepo) FindByHadistID(hadistID uint, page, limit int) ([]domain.Doa, int64, error) {
	var doas []domain.Doa
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&domain.Doa{}).Where("hadist_id = ?", hadistID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Hadist").Limit(limit).Offset(offset).Order("created_at desc").Find(&doas).Error; err != nil {
		return nil, 0, err
	}

	return doas, total, nil
}

func (r *DoaRepo) FindRandom() (*domain.Doa, error) {
	var doa domain.Doa
	if err := r.db.Preload("Hadist").Order("RANDOM()").First(&doa).Error; err != nil {
		return nil, err
	}
	return &doa, nil
}

func (r *DoaRepo) Update(doa *domain.Doa) error {
	return r.db.Save(doa).Error
}

func (r *DoaRepo) Delete(id uint) error {
	return r.db.Delete(&domain.Doa{}, id).Error
}

func (r *DoaRepo) IncrementListeningCount(id uint) error {
	return r.db.Model(&domain.Doa{}).Where("id = ?", id).UpdateColumn("listening_count", gorm.Expr("listening_count + ?", 1)).Error
}

// Like Operations
func (r *DoaRepo) CreateLike(like *domain.DoaLike) error {
	return r.db.Create(like).Error
}

func (r *DoaRepo) FindLikeByUserAndDoa(userID, doaID uint) (*domain.DoaLike, error) {
	var like domain.DoaLike
	err := r.db.Where("user_id = ? AND doa_id = ?", userID, doaID).First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}

func (r *DoaRepo) DeleteLike(id uint) error {
	return r.db.Delete(&domain.DoaLike{}, id).Error
}

func (r *DoaRepo) IncrementLikeCount(doaID uint) error {
	return r.db.Model(&domain.Doa{}).Where("id = ?", doaID).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
}

func (r *DoaRepo) DecrementLikeCount(doaID uint) error {
	return r.db.Model(&domain.Doa{}).Where("id = ?", doaID).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error
}

// Bookmark Operations
func (r *DoaRepo) CreateBookmark(bookmark *domain.DoaBookmark) error {
	return r.db.Create(bookmark).Error
}

func (r *DoaRepo) FindBookmarkByUserAndDoa(userID, doaID uint) (*domain.DoaBookmark, error) {
	var bookmark domain.DoaBookmark
	err := r.db.Where("user_id = ? AND doa_id = ?", userID, doaID).First(&bookmark).Error
	if err != nil {
		return nil, err
	}
	return &bookmark, nil
}

func (r *DoaRepo) DeleteBookmark(id uint) error {
	return r.db.Delete(&domain.DoaBookmark{}, id).Error
}

func (r *DoaRepo) IncrementBookmarkCount(doaID uint) error {
	return r.db.Model(&domain.Doa{}).Where("id = ?", doaID).UpdateColumn("bookmark_count", gorm.Expr("bookmark_count + ?", 1)).Error
}

func (r *DoaRepo) DecrementBookmarkCount(doaID uint) error {
	return r.db.Model(&domain.Doa{}).Where("id = ?", doaID).UpdateColumn("bookmark_count", gorm.Expr("bookmark_count - ?", 1)).Error
}

func (r *DoaRepo) GetUserBookmarks(userID uint, page, limit int) ([]domain.DoaBookmark, int64, error) {
	var bookmarks []domain.DoaBookmark
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&domain.DoaBookmark{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Doa").Preload("Doa.Hadist").
		Where("user_id = ?", userID).
		Limit(limit).Offset(offset).
		Order("created_at desc").
		Find(&bookmarks).Error

	if err != nil {
		return nil, 0, err
	}

	return bookmarks, total, nil
}
