package service

import (
	"errors"

	"github.com/project/wayt/internal/model"
	"github.com/project/wayt/internal/repository"
)

type CounterService interface {
	Create(branchID uint, name, prefix string) (*model.Counter, error)
	ListByBranch(branchID uint) ([]model.Counter, error)
	FindByID(id uint) (*model.Counter, error)
	Update(id uint, name, prefix string, isActive bool) (*model.Counter, error)
	Delete(id uint) error
}

type counterService struct {
	repo       repository.CounterRepository
	branchRepo repository.BranchRepository
}

func NewCounterService(repo repository.CounterRepository, branchRepo repository.BranchRepository) CounterService {
	return &counterService{repo: repo, branchRepo: branchRepo}
}

func (s *counterService) Create(branchID uint, name, prefix string) (*model.Counter, error) {
	if name == "" || prefix == "" {
		return nil, errors.New("name and prefix are required")
	}
	if _, err := s.branchRepo.FindByID(branchID); err != nil {
		return nil, errors.New("branch not found")
	}
	counter := &model.Counter{
		BranchID: branchID,
		Name:     name,
		Prefix:   prefix,
		IsActive: true,
	}
	if err := s.repo.Create(counter); err != nil {
		return nil, err
	}
	return counter, nil
}

func (s *counterService) ListByBranch(branchID uint) ([]model.Counter, error) {
	return s.repo.FindByBranch(branchID)
}

func (s *counterService) FindByID(id uint) (*model.Counter, error) {
	counter, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("counter not found")
	}
	return counter, nil
}

func (s *counterService) Update(id uint, name, prefix string, isActive bool) (*model.Counter, error) {
	counter, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("counter not found")
	}
	if name != "" {
		counter.Name = name
	}
	if prefix != "" {
		counter.Prefix = prefix
	}
	counter.IsActive = isActive
	if err := s.repo.Update(counter); err != nil {
		return nil, err
	}
	return counter, nil
}

func (s *counterService) Delete(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return errors.New("counter not found")
	}
	return s.repo.Delete(id)
}
