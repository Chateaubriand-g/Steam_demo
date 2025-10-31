package services

import (
	"steam-backend/models"
	"steam-backend/repositories"
)

type WishlistService interface {
	AddToWishlist(userID, appID uint64) error
	RemoveFromWishlist(userID, appID uint64) error
	GetWishlistSize(userID uint64) (int64, error)
	GetWishlist(userID uint64) ([]models.WishlistItemDto, error)
	IsInWishlist(userID, appID uint64) (bool, error)
	SortWishlist(userID uint64, sortItems []models.SortItem) error
}

type wishlistService struct {
	wishlistRepo repositories.WishlistRepository
	appRepo      repositories.AppRepository
}

func NewWishlistService(wishrepo repositories.WishlistRepository, apprepo repositories.AppRepository) WishlistService {
	return &wishlistService{
		wishlistRepo: wishrepo,
		appRepo:      apprepo,
	}
}

func (s *wishlistService) AddToWishlist(userID, appID uint64) error {
	_, err := s.appRepo.FindByID(appID)
	if err != nil {
		return err
	}

	return s.wishlistRepo.AddItem(userID, appID)
}

func (s *wishlistService) RemoveFromWishlist(userID, appID uint64) error {
	_, err := s.appRepo.FindByID(appID)
	if err != nil {
		return err
	}

	return s.wishlistRepo.RemoveItem(userID, appID)
}

func (s *wishlistService) GetWishlistSize(userID uint64) (int64, error) {
	return s.wishlistRepo.GetItemCount(userID)
}

func (s *wishlistService) GetWishlist(userID uint64) ([]models.WishlistItemDto, error) {
	wishlist, err := s.wishlistRepo.GetWishlist(userID)
	if err != nil {
		return nil, err
	}
	wishlist_dto := s.convertToDtos(wishlist)
	return wishlist_dto, nil
}

func (s *wishlistService) IsInWishlist(userID, appID uint64) (bool, error) {
	return s.wishlistRepo.IsInWishList(userID, appID)
}

func (s *wishlistService) SortWishlist(userID uint64, sortItems []models.SortItem) error {
	return s.wishlistRepo.UpdateItemOrder(userID, sortItems)
}

func (s *wishlistService) convertToDtos(items []models.WishlistItem) []models.WishlistItemDto {
	var res []models.WishlistItemDto
	for _, item := range items {
		dto, err := s.convertToDto(item)
		if err != nil {
			continue
		}
		res = append(res, dto)
	}
	return res
}

func (s *wishlistService) convertToDto(item models.WishlistItem) (models.WishlistItemDto, error) {
	appID := item.AppID

	app, err := s.appRepo.FindByID(appID)
	if err != nil {
		return models.WishlistItemDto{}, err
	}

	dto := models.WishlistItemDto{
		AppID:        appID,
		Name:         app.Name,
		Price:        app.Price,
		Discount:     app.Discount,
		ImageURL:     app.ImageURL,
		PositiveRate: app.PositiveRate,
		SortOrder:    item.SortOrder,
	}

	return dto, nil
}
