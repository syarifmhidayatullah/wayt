package service

import (
	"errors"

	"github.com/project/wayt/internal/model"
	"github.com/project/wayt/internal/repository"
)

type BranchService interface {
	Create(name, prefix string) (*model.Branch, error)
	List() ([]model.Branch, error)
	Update(id uint, name, prefix string, isActive bool) (*model.Branch, error)
	Delete(id uint) error
}

type branchService struct {
	repo repository.BranchRepository
}

func NewBranchService(repo repository.BranchRepository) BranchService {
	return &branchService{repo: repo}
}

func (s *branchService) Create(name, prefix string) (*model.Branch, error) {
	if name == "" || prefix == "" {
		return nil, errors.New("name and prefix are required")
	}
	branch := &model.Branch{
		Name:     name,
		Prefix:   prefix,
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

func (s *branchService) Update(id uint, name, prefix string, isActive bool) (*model.Branch, error) {
	branch, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("branch not found")
	}
	if name != "" {
		branch.Name = name
	}
	if prefix != "" {
		branch.Prefix = prefix
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
