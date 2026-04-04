package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/project/wayt/internal/model"
	"github.com/project/wayt/internal/repository"
)

type RegisterResult struct {
	QueueID     uint   `json:"queue_id"`
	QueueNumber string `json:"queue_number"`
	BranchName  string `json:"branch_name"`
	Position    int64  `json:"position"`
	PeopleAhead int64  `json:"people_ahead"`
}

type QueueStatusResult struct {
	QueueNumber    string            `json:"queue_number"`
	Status         model.QueueStatus `json:"status"`
	CurrentServing string            `json:"current_serving"`
	PeopleAhead    int64             `json:"people_ahead"`
}

type QueueService interface {
	Register(qrToken string) (*RegisterResult, error)
	ScanRegister(qrToken string) (*RegisterResult, error)
	StatusByID(queueID uint) (*QueueStatusResult, error)
	Status(qrToken string) (*QueueStatusResult, error)
	CallNext(branchID uint) (*model.Queue, error)
	ListByBranch(branchID uint) ([]model.Queue, error)
	Reset(branchID uint) error
}

type queueService struct {
	queueRepo  repository.QueueRepository
	qrRepo     repository.QRCodeRepository
	branchRepo repository.BranchRepository
}

func NewQueueService(queueRepo repository.QueueRepository, qrRepo repository.QRCodeRepository, branchRepo repository.BranchRepository) QueueService {
	return &queueService{queueRepo: queueRepo, qrRepo: qrRepo, branchRepo: branchRepo}
}

func (s *queueService) ScanRegister(qrToken string) (*RegisterResult, error) {
	return s.Register(qrToken)
}

func (s *queueService) Register(qrToken string) (*RegisterResult, error) {
	qr, err := s.qrRepo.FindByToken(qrToken)
	if err != nil {
		return nil, errors.New("QR code not found")
	}
	if !qr.IsActive {
		return nil, errors.New("QR code is no longer active")
	}
	if time.Now().After(qr.ExpiredAt) {
		return nil, errors.New("QR code has expired")
	}

	branch, err := s.branchRepo.FindByID(qr.BranchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	branch.LastNumber++
	queueNumber := fmt.Sprintf("%s-%03d", branch.Prefix, branch.LastNumber)

	if err := s.branchRepo.Update(branch); err != nil {
		return nil, err
	}

	queue := &model.Queue{
		BranchID:    branch.ID,
		QRToken:     qrToken,
		QueueNumber: queueNumber,
		Status:      model.QueueStatusWaiting,
	}
	if err := s.queueRepo.Create(queue); err != nil {
		return nil, err
	}

	ahead, err := s.queueRepo.CountWaitingAhead(branch.ID, queue.ID)
	if err != nil {
		ahead = 0
	}

	return &RegisterResult{
		QueueID:     queue.ID,
		QueueNumber: queueNumber,
		BranchName:  branch.Name,
		Position:    ahead + 1,
		PeopleAhead: ahead,
	}, nil
}

func (s *queueService) StatusByID(queueID uint) (*QueueStatusResult, error) {
	queue, err := s.queueRepo.FindByID(queueID)
	if err != nil {
		return nil, errors.New("queue not found")
	}

	branch, err := s.branchRepo.FindByID(queue.BranchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	currentServing := fmt.Sprintf("%s-%03d", branch.Prefix, branch.CurrentNumber)

	var ahead int64
	if queue.Status == model.QueueStatusWaiting {
		ahead, _ = s.queueRepo.CountWaitingAhead(branch.ID, queue.ID)
	}

	return &QueueStatusResult{
		QueueNumber:    queue.QueueNumber,
		Status:         queue.Status,
		CurrentServing: currentServing,
		PeopleAhead:    ahead,
	}, nil
}

func (s *queueService) Status(qrToken string) (*QueueStatusResult, error) {
	queue, err := s.queueRepo.FindByQRToken(qrToken)
	if err != nil {
		return nil, errors.New("queue not found for this token")
	}

	branch, err := s.branchRepo.FindByID(queue.BranchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	currentServing := fmt.Sprintf("%s-%03d", branch.Prefix, branch.CurrentNumber)

	var ahead int64
	if queue.Status == model.QueueStatusWaiting {
		ahead, _ = s.queueRepo.CountWaitingAhead(branch.ID, queue.ID)
	}

	return &QueueStatusResult{
		QueueNumber:    queue.QueueNumber,
		Status:         queue.Status,
		CurrentServing: currentServing,
		PeopleAhead:    ahead,
	}, nil
}

func (s *queueService) CallNext(branchID uint) (*model.Queue, error) {
	if _, err := s.branchRepo.FindByID(branchID); err != nil {
		return nil, errors.New("branch not found")
	}

	next, err := s.queueRepo.FindNextWaiting(branchID)
	if err != nil {
		return nil, errors.New("no waiting queue found")
	}

	if err := s.queueRepo.UpdateStatus(next.ID, model.QueueStatusCalled); err != nil {
		return nil, err
	}
	if err := s.branchRepo.IncrementCurrentNumber(branchID); err != nil {
		return nil, err
	}

	next.Status = model.QueueStatusCalled
	return next, nil
}

func (s *queueService) ListByBranch(branchID uint) ([]model.Queue, error) {
	if _, err := s.branchRepo.FindByID(branchID); err != nil {
		return nil, errors.New("branch not found")
	}
	return s.queueRepo.FindActiveByBranch(branchID)
}

func (s *queueService) Reset(branchID uint) error {
	if _, err := s.branchRepo.FindByID(branchID); err != nil {
		return errors.New("branch not found")
	}
	if err := s.queueRepo.ExpireByBranch(branchID); err != nil {
		return err
	}
	if err := s.qrRepo.DeactivateByBranch(branchID); err != nil {
		return err
	}
	return s.branchRepo.ResetNumbers(branchID)
}
