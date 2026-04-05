package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/project/wayt/internal/model"
	"github.com/project/wayt/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(name, phone, password string) (*model.User, error)
	Login(phone, password string) (string, error)
	FindByID(id uint) (*model.User, error)
}

type userService struct {
	repo      repository.UserRepository
	jwtSecret []byte
}

func NewUserService(repo repository.UserRepository, jwtSecret string) UserService {
	return &userService{repo: repo, jwtSecret: []byte(jwtSecret)}
}

func (s *userService) Register(name, phone, password string) (*model.User, error) {
	if name == "" || phone == "" || password == "" {
		return nil, errors.New("name, phone, dan password wajib diisi")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Name:     name,
		Phone:    phone,
		Password: string(hashed),
	}
	if err := s.repo.Create(user); err != nil {
		return nil, errors.New("nomor telepon sudah terdaftar")
	}
	return user, nil
}

func (s *userService) Login(phone, password string) (string, error) {
	user, err := s.repo.FindByPhone(phone)
	if err != nil {
		return "", errors.New("nomor telepon atau password salah")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("nomor telepon atau password salah")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"name":  user.Name,
		"phone": user.Phone,
		"type":  "user",
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", errors.New("gagal membuat token")
	}
	return signed, nil
}

func (s *userService) FindByID(id uint) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return user, nil
}
