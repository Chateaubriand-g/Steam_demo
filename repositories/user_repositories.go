package repositories

import (
	"errors"
	"steam-backend/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint64) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint64) error
	SearchUsers(keyword string, limit int) ([]models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) (UserRepository, error) {
	if db == nil {
		return nil, errors.New("db to userRepository is nil")
	}
	return &userRepository{db: db}, nil
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint64) (*models.User, error) {
	var res models.User

	query := r.db.Where("userId = ?", id)
	err := query.Find(&res).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var res models.User

	query := r.db.Where("email = ?", email)
	err := query.Find(&res).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var res models.User

	query := r.db.Where("userName = ?", username)
	err := query.Find(&res).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint64) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) SearchUsers(keyword string, limit int) ([]models.User, error) {
	var res []models.User

	query := r.db.Where("userName like ?", keyword+"%")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&res).Error
	return res, err
}
