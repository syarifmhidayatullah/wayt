package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/project/wayt/internal/model"
	"github.com/project/wayt/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(username, password string) (string, error)
	SeedAdmin(username, password string) error
	ListUsers() ([]model.AdminUser, error)
	CreateUser(username, password string, role model.AdminRole, branchID *uint) (*model.AdminUser, error)
	UpdateUser(id uint, username string, role model.AdminRole, password string, branchID *uint) (*model.AdminUser, error)
	DeleteUser(id uint, requesterID uint) error
}

type authService struct {
	repo      repository.AdminUserRepository
	jwtSecret []byte
}

func NewAuthService(repo repository.AdminUserRepository, jwtSecret string) AuthService {
	return &authService{repo: repo, jwtSecret: []byte(jwtSecret)}
}

func (s *authService) Login(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", errors.New("username atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("username atau password salah")
	}

	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"role":     string(user.Role),
		"exp":      time.Now().Add(8 * time.Hour).Unix(),
	}
	if user.BranchID != nil {
		claims["branch_id"] = *user.BranchID
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", errors.New("gagal membuat token")
	}

	return signed, nil
}

func (s *authService) SeedAdmin(username, password string) error {
	exists, err := s.repo.ExistsAny()
	if err != nil || exists {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.Create(&model.AdminUser{
		Username: username,
		Role:     model.RoleSuperAdmin,
		Password: string(hashed),
	})
}

func (s *authService) ListUsers() ([]model.AdminUser, error) {
	return s.repo.FindAll()
}

func (s *authService) CreateUser(username, password string, role model.AdminRole, branchID *uint) (*model.AdminUser, error) {
	if username == "" || password == "" {
		return nil, errors.New("username dan password wajib diisi")
	}
	if role != model.RoleSuperAdmin && role != model.RoleAdmin {
		role = model.RoleAdmin
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.AdminUser{
		Username: username,
		Role:     role,
		Password: string(hashed),
		BranchID: branchID,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, errors.New("username sudah digunakan")
	}
	return user, nil
}

func (s *authService) UpdateUser(id uint, username string, role model.AdminRole, password string, branchID *uint) (*model.AdminUser, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if username != "" {
		user.Username = username
	}
	if role == model.RoleSuperAdmin || role == model.RoleAdmin {
		user.Role = role
	}
	if password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashed)
	}
	// branchID nil means "no change"; use a sentinel pointer-of-zero to clear
	if branchID != nil {
		if *branchID == 0 {
			user.BranchID = nil
		} else {
			user.BranchID = branchID
		}
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) DeleteUser(id uint, requesterID uint) error {
	if id == requesterID {
		return errors.New("tidak bisa menghapus akun sendiri")
	}
	if _, err := s.repo.FindByID(id); err != nil {
		return errors.New("user tidak ditemukan")
	}
	return s.repo.Delete(id)
}
