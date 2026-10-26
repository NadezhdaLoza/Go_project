package implementation

import (
    "github.com/nadezhda/documents/services/base"
    "github.com/nadezhda/documents/services/work/backend/repo"
    svc "github.com/nadezhda/documents/services/work/backend/service"
    "github.com/nadezhda/documents/services/work/domain"
)

type service struct {
    repo repo.WorkRepo
}

func NewService(repo repo.WorkRepo) svc.WorkService {
    return &service{
        repo: repo,
    }
}

func (s *service) CreateWork(name string, rate float64, unit string) (*domain.Work, error) {
    work := domain.Work{
        Name: name,
        Rate: rate,
        Unit: unit,
    }
    result, err := s.repo.CreateWork(work)
    base.LogError(err)
    return result, err
}

func (s *service) GetWorks() ([]domain.Work, error) {
    result, err := s.repo.GetWorks()
    base.LogError(err)
    return result, err
}

func (s *service) UpdateWork(work domain.Work) (*domain.Work, error) {
    result, err := s.repo.UpdateWork(work)
    base.LogError(err)
    return result, err
}

func (s *service) DeleteWork(id int) error {
    err := s.repo.DeleteWork(id)
    base.LogError(err)
    return err
}

func (s *service) WorkById(id int) (*domain.Work, error) {
    work, err := s.repo.WorkById(id)
    base.LogError(err)
    return work, err
}