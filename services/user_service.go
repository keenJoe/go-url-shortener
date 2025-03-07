package services

import (
	"errors"

	"github.com/keenJoe/go-url-shortener/models"
	"github.com/keenJoe/go-url-shortener/pkg/database"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserServiceInterface 用户服务接口
type UserServiceInterface interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id string) (*models.User, error)
	CreateUser(req models.CreateUserRequest) (*models.User, error)
	UpdateUser(id string, req models.UpdateUserRequest) (*models.User, error)
	DeleteUser(id string) error
}

// UserService 用户服务实现
type UserService struct {
	db *gorm.DB
}

func NewUserService() UserServiceInterface {
	return &UserService{
		db: database.DB,
	}
}

// 实现 UserServiceInterface 接口的所有方法
func (s *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := s.db.Find(&users)
	return users, result.Error
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	var user models.User
	result := s.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, result.Error
	}
	return &user, nil
}

func (s *UserService) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	// 检查用户名是否已存在
	var existingUser models.User
	if result := s.db.Where("username = ?", req.Username).First(&existingUser); result.Error == nil {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	result := s.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (s *UserService) UpdateUser(id string, req models.UpdateUserRequest) (*models.User, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if req.Username != "" {
		// 检查新用户名是否与其他用户冲突
		var existingUser models.User
		if result := s.db.Where("username = ? AND id != ?", req.Username, id).First(&existingUser); result.Error == nil {
			return nil, errors.New("用户名已存在")
		}
		updates["username"] = req.Username
	}

	if req.Email != "" {
		updates["email"] = req.Email
	}

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		updates["password"] = string(hashedPassword)
	}

	result := s.db.Model(user).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (s *UserService) DeleteUser(id string) error {
	result := s.db.Delete(&models.User{}, id)
	return result.Error
}
