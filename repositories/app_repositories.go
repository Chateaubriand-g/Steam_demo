package repositories

import (
	"errors"
	"steam-backend/models"

	"gorm.io/gorm"
)

type AppRepository interface {
	FindByID(id uint64) (*models.App, error)
	FindRecommendations(limit int) ([]models.App, error)
	FindSpecials(limit int) ([]models.App, error)
	SearchSuggestions(key string, limit int) ([]models.App, error)
	SearchApps(key string, page, pageSize int) ([]models.App, int64, error)
}

type appRepository struct {
	db *gorm.DB
}

func NewAppRepository(db *gorm.DB) (AppRepository, error) {
	if db == nil {
		return nil, errors.New("gorm.DB is nil")
	}
	return &appRepository{db: db}, nil
}

func (r *appRepository) FindByID(id uint64) (*models.App, error) {
	var res models.App
	// ("appId = ?",id)会将id转换为字符，避免SQL注入攻击
	err := r.db.Where("appId = ?", id).First(&res).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *appRepository) FindRecommendations(limit int) ([]models.App, error) {
	var res []models.App

	query := r.db.Order("positiveRate DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	//query接收的是gorm的查询构造器，用于积累查询条件，当执行Find,First,Scan等执行方法时
	//query会根据积累的查询条件生成SQL语句并执行，将结果映射到传入的结构体中
	err := query.Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *appRepository) FindSpecials(limit int) ([]models.App, error) {
	var res []models.App

	query := r.db.Where("discount > 0").Order("discount DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *appRepository) SearchSuggestions(key string, limit int) ([]models.App, error) {
	var res []models.App

	query := r.db.Where("name like ?", key).Order("positiveRate DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *appRepository) SearchApps(key string, page, pageSize int) ([]models.App, int64, error) {
	var res []models.App
	var resCounts int64

	query := r.db.Model(&models.App{}).Where("name like ?", key+"%")
	err := query.Count(&resCounts).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query = query.Order("positiveRare DESC").Limit(pageSize).Offset(offset)
	err = query.Find(&res).Error
	if err != nil {
		return nil, 0, err
	}
	return res, resCounts, nil
}
