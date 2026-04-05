package service

import (
	"errors"

	"github.com/project/wayt/internal/model"
	"github.com/project/wayt/internal/repository"
)

type BranchService interface {
	Create(name string) (*model.Branch, error)
	List() ([]model.Branch, error)
	FindByID(id uint) (*model.Branch, error)
	Update(id uint, name string, isActive bool) (*model.Branch, error)
	Delete(id uint) error
}

type branchService struct {
	repo repository.BranchRepository
}

func NewBranchService(repo repository.BranchRepository) BranchService {
	return &branchService{repo: repo}
}

func (s *branchService) Create(name string) (*model.Branch, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	branch := &model.Branch{
		Name:     name,
		IsActive: true,
	}
	if err := s.repo.Create(branch); err != nil {
		return nil, err
	}
	return branch, nil
}

func (s *branchService) List() ([]model.Branch, error) {
	return s.repo.FindAll()
}

func (s *branchService) FindByID(id uint) (*model.Branch, error) {
	branch, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("branch not found")
	}
	return branch, nil
}

func (s *branchService) Update(id uint, name string, isActive bool) (*model.Branch, error) {
	branch, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("branch not found")
	}
	if name != "" {
		branch.Name = name
	}
	branch.IsActive = isActive
	if err := s.repo.Update(branch); err != nil {
		return nil, err
	}
	return branch, nil
}

func (s *branchService) Delete(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return errors.New("branch not found")
	}
	return s.repo.Delete(id)
}
