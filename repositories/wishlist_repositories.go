package repositories

import (
	"errors"
	"steam-backend/models"

	"gorm.io/gorm"
)

type WishlistRepository interface {
	AddItem(userID, appID uint64) error
	RemoveItem(userID, appID uint64) error
	GetItemCount(userID uint64) (int64, error)
	GetWishlist(userID uint64) ([]models.WishlistItem, error)
	IsInWishList(userID, appID uint64) (bool, error)
	UpdateItemOrder(userID uint64, sortItems []models.SortItem) error
}

type wishlistRepository struct {
	db *gorm.DB
}

func NewWishlistRepository(db *gorm.DB) (WishlistRepository, error) {
	if db == nil {
		return nil, errors.New("db to wishlistRepository is nil")
	}
	return &wishlistRepository{db: db}, nil
}

func (r *wishlistRepository) AddItem(userID, appID uint64) error {
	exists, err := r.IsInWishList(userID, appID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	var maxOrder int64
	r.db.Model(&models.WishlistItem{}).Where("userId = ?", userID).Select(
		"COALESCE(MAX(sortOrder),-1)").Scan(&maxOrder)

	newItem := models.WishlistItem{
		UserID:    userID,
		AppID:     appID,
		SortOrder: maxOrder + 1,
	}

	return r.db.Create(&newItem).Error
}

func (r *wishlistRepository) RemoveItem(userID, appID uint64) error {
	return r.db.Where("userId = ? and appId = ?", userID, appID).Delete(&models.WishlistItem{}).Error
}

func (r *wishlistRepository) GetItemCount(userID uint64) (int64, error) {
	var count int64
	err := r.db.Where("userId = ?", userID).Count(&count).Error
	return count, err
}

func (r *wishlistRepository) GetWishlist(userID uint64) ([]models.WishlistItem, error) {
	var res []models.WishlistItem

	query := r.db.Where("userId = ?", userID)
	err := query.Find(&res).Error
	return res, err
}

func (r *wishlistRepository) IsInWishList(userID, appID uint64) (bool, error) {
	var count int64

	err := r.db.Where("userId = ? and appId = ?", userID, appID).Count(&count).Error

	if count > 0 {
		return true, err
	}
	return false, err
}

func (r *wishlistRepository) UpdateItemOrder(userID uint64, sortItems []models.SortItem) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, item := range sortItems {
		if err := tx.Model(&models.WishlistItem{}).Where("userId = ? and appId = ?", userID, item.AppID).Update(
			"sortOrder", item.SortOrder).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
