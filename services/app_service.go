package services

import (
	"steam-backend/models"
	"steam-backend/repositories"
)

type AppService interface {
	GetRecommendations(limit int) ([]models.AppDto, error)
	GetSpecials(limit int) ([]models.AppDto, error)
	GetSearchSuggestions(keyword string, limit int) ([]models.AppDto, error)
	GetAppByID(id uint64) (*models.AppDto, error)
}

type appService struct {
	appRepo repositories.AppRepository
}

func NewAPPService(appRepo repositories.AppRepository) AppService {
	return &appService{appRepo: appRepo}
}

func (s *appService) GetRecommendations(limit int) ([]models.AppDto, error) {
	res, err := s.appRepo.FindRecommendations(limit)
	if err != nil {
		return nil, err
	}
	return s.convertToAppDtos(res), nil
}

func (s *appService) GetSpecials(limit int) ([]models.AppDto, error) {
	res, err := s.appRepo.FindSpecials(limit)
	if err != nil {
		return nil, err
	}
	return s.convertToAppDtos(res), nil
}

func (s *appService) GetSearchSuggestions(keyword string, limit int) ([]models.AppDto, error) {
	res, err := s.appRepo.SearchSuggestions(keyword, limit)
	if err != nil {
		return nil, err
	}
	return s.convertToAppDtos(res), nil
}

func (s *appService) GetAppByID(id uint64) (*models.AppDto, error) {
	res, err := s.appRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	appDto := s.convertToAppDto(res)
	return &appDto, nil
}

func (s *appService) convertToAppDtos(apps []models.App) []models.AppDto {
	res := make([]models.AppDto, len(apps))
	for i, app := range apps {
		res[i] = s.convertToAppDto(&app)
	}
	return res
}

func (s *appService) convertToAppDto(app *models.App) models.AppDto {
	return models.AppDto{
		AppId:        app.AppId,
		Name:         app.Name,
		Price:        app.Price,
		Discount:     app.Discount,
		ImageURL:     app.ImageURL,
		PositiveRate: app.PositiveRate,
	}
}
