package services

import (
	"errors"
	"steam-backend/config"
	"steam-backend/models"
	"steam-backend/repositories"
	"steam-backend/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Register(userDTO *models.JoinRequestDto) (*models.User, error)
	Login(loginDTO *models.LoginRequestDto) (string, *models.User, error)
	GetUserInfo(userID uint64) (*models.User, error)
	ChechUserNameAvailable(username string) (bool, error)
	SearchUsers(keyword string) ([]models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
	config   config.Config
}

func NewUserService(repo repositories.UserRepository, conf config.Config) UserService {
	return &userService{
		userRepo: repo,
		config:   conf,
	}
}

func (s *userService) Register(joinRequestDTO *models.JoinRequestDto) (*models.User, error) {
	userName, _ := s.userRepo.FindByUsername(joinRequestDTO.UserName)
	if userName != nil {
		return nil, errors.New("userName been used")
	}
	email, _ := s.userRepo.FindByEmail(joinRequestDTO.Email)
	if email != nil {
		return nil, errors.New("email been used")
	}
	hashPassword, err := HashPassword(joinRequestDTO.PassWord)
	if err != nil {
		return nil, errors.New("hashpassword failed")
	}
	newUser := &models.User{
		Email:    joinRequestDTO.Email,
		UserName: hashPassword,
		PassWord: joinRequestDTO.PassWord,
	}
	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}

func (s *userService) Login(loginDTO *models.LoginRequestDto) (string, *models.User, error) {
	user, _ := s.userRepo.FindByUsername(loginDTO.UserName)
	if user == nil {
		return "", nil, errors.New("username is not exists")
	}

	compare := CheckPassword(loginDTO.Password, user.PassWord)
	if !compare {
		return "", nil, errors.New("password is incorrect")
	}

	token, err := utils.GenerateToken(user.UserID, s.config.JWTSecret)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *userService) GetUserInfo(userID uint64) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *userService) ChechUserNameAvailable(username string) (bool, error) {
	_, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}
		return false, err
	}
	return false, errors.New("userName exists")
}

func (s *userService) SearchUsers(keyword string) ([]models.User, error) {
	return s.userRepo.SearchUsers(keyword, 20)
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
