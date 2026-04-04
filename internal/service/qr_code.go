package service

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/project/wayt/config"
	"github.com/project/wayt/internal/model"
	"github.com/project/wayt/internal/repository"
	qrcode "github.com/skip2/go-qrcode"
)

type QRCodeResult struct {
	Token      string    `json:"token"`
	QRImageURL string    `json:"qr_image_url"`
	ExpiredAt  time.Time `json:"expired_at"`
}

type QRCodeService interface {
	Generate(branchID uint) (*QRCodeResult, error)
}

type qrCodeService struct {
	qrRepo     repository.QRCodeRepository
	branchRepo repository.BranchRepository
	cfg        config.QRConfig
}

func NewQRCodeService(qrRepo repository.QRCodeRepository, branchRepo repository.BranchRepository, cfg config.QRConfig) QRCodeService {
	return &qrCodeService{qrRepo: qrRepo, branchRepo: branchRepo, cfg: cfg}
}

func (s *qrCodeService) Generate(branchID uint) (*QRCodeResult, error) {
	branch, err := s.branchRepo.FindByID(branchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}
	if !branch.IsActive {
		return nil, errors.New("branch is not active")
	}

	token := uuid.New().String()
	expiredAt := time.Now().Add(time.Duration(s.cfg.ExpiredHours) * time.Hour)

	filename := fmt.Sprintf("%s.png", token)
	filePath := filepath.Join(s.cfg.StoragePath, filename)

	// QR content is a URL; scanning it directly registers the queue via GET /q/:token
	qrContent := fmt.Sprintf("%s/q/%s", s.cfg.PublicBaseURL, token)
	if err := qrcode.WriteFile(qrContent, qrcode.Medium, 256, filePath); err != nil {
		return nil, fmt.Errorf("failed to generate QR image: %w", err)
	}

	qr := &model.QRCode{
		BranchID:  branchID,
		Token:     token,
		IsActive:  true,
		ExpiredAt: expiredAt,
	}
	if err := s.qrRepo.Create(qr); err != nil {
		return nil, err
	}

	return &QRCodeResult{
		Token:      token,
		QRImageURL: fmt.Sprintf("%s/%s", s.cfg.BaseURL, filename),
		ExpiredAt:  expiredAt,
	}, nil
}
