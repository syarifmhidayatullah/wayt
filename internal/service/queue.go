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
	CounterName string `json:"counter_name"`
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
	CallNext(counterID uint) (*model.Queue, error)
	ListByCounter(counterID uint) ([]model.Queue, error)
	Reset(counterID uint) error
}

type queueService struct {
	queueRepo   repository.QueueRepository
	qrRepo      repository.QRCodeRepository
	counterRepo repository.CounterRepository
}

func NewQueueService(queueRepo repository.QueueRepository, qrRepo repository.QRCodeRepository, counterRepo repository.CounterRepository) QueueService {
	return &queueService{queueRepo: queueRepo, qrRepo: qrRepo, counterRepo: counterRepo}
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

	counter, err := s.counterRepo.FindByID(qr.CounterID)
	if err != nil {
		return nil, errors.New("counter not found")
	}

	counter.LastNumber++
	queueNumber := fmt.Sprintf("%s-%03d", counter.Prefix, counter.LastNumber)

	if err := s.counterRepo.Update(counter); err != nil {
		return nil, err
	}

	queue := &model.Queue{
		BranchID:    counter.BranchID,
		CounterID:   counter.ID,
		QRToken:     qrToken,
		QueueNumber: queueNumber,
		Status:      model.QueueStatusWaiting,
	}
	if err := s.queueRepo.Create(queue); err != nil {
		return nil, err
	}

	ahead, err := s.queueRepo.CountWaitingAhead(counter.ID, queue.ID)
	if err != nil {
		ahead = 0
	}

	branchName := ""
	if qr.Branch != nil {
		branchName = qr.Branch.Name
	}

	return &RegisterResult{
		QueueID:     queue.ID,
		QueueNumber: queueNumber,
		BranchName:  branchName,
		CounterName: counter.Name,
		Position:    ahead + 1,
		PeopleAhead: ahead,
	}, nil
}

func (s *queueService) StatusByID(queueID uint) (*QueueStatusResult, error) {
	queue, err := s.queueRepo.FindByID(queueID)
	if err != nil {
		return nil, errors.New("queue not found")
	}

	counter, err := s.counterRepo.FindByID(queue.CounterID)
	if err != nil {
		return nil, errors.New("counter not found")
	}

	currentServing := fmt.Sprintf("%s-%03d", counter.Prefix, counter.CurrentNumber)

	var ahead int64
	if queue.Status == model.QueueStatusWaiting {
		ahead, _ = s.queueRepo.CountWaitingAhead(counter.ID, queue.ID)
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

	counter, err := s.counterRepo.FindByID(queue.CounterID)
	if err != nil {
		return nil, errors.New("counter not found")
	}

	currentServing := fmt.Sprintf("%s-%03d", counter.Prefix, counter.CurrentNumber)

	var ahead int64
	if queue.Status == model.QueueStatusWaiting {
		ahead, _ = s.queueRepo.CountWaitingAhead(counter.ID, queue.ID)
	}

	return &QueueStatusResult{
		QueueNumber:    queue.QueueNumber,
		Status:         queue.Status,
		CurrentServing: currentServing,
		PeopleAhead:    ahead,
	}, nil
}

func (s *queueService) CallNext(counterID uint) (*model.Queue, error) {
	if _, err := s.counterRepo.FindByID(counterID); err != nil {
		return nil, errors.New("counter not found")
	}

	next, err := s.queueRepo.FindNextWaiting(counterID)
	if err != nil {
		return nil, errors.New("no waiting queue found")
	}

	if err := s.queueRepo.UpdateStatus(next.ID, model.QueueStatusCalled); err != nil {
		return nil, err
	}
	if err := s.counterRepo.IncrementCurrentNumber(counterID); err != nil {
		return nil, err
	}

	next.Status = model.QueueStatusCalled
	return next, nil
}

func (s *queueService) ListByCounter(counterID uint) ([]model.Queue, error) {
	if _, err := s.counterRepo.FindByID(counterID); err != nil {
		return nil, errors.New("counter not found")
	}
	return s.queueRepo.FindActiveByCounter(counterID)
}

func (s *queueService) Reset(counterID uint) error {
	if _, err := s.counterRepo.FindByID(counterID); err != nil {
		return errors.New("counter not found")
	}
	if err := s.queueRepo.ExpireByCounter(counterID); err != nil {
		return err
	}
	if err := s.qrRepo.DeactivateByCounter(counterID); err != nil {
		return err
	}
	return s.counterRepo.ResetNumbers(counterID)
}
