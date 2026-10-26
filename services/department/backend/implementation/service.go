package implementation

import (
    "github.com/nadezhda/documents/services/base"
    "github.com/nadezhda/documents/services/department/backend/repo"
    svc "github.com/nadezhda/documents/services/department/backend/service"
    "github.com/nadezhda/documents/services/department/domain"
)

type service struct {
    repo repo.DepartmentRepo
}

func NewService(repo repo.DepartmentRepo) svc.DepartmentService {
    return &service{
        repo: repo,
    }
}

func (s *service) CreateDepartment(name, workConditionsClass, workDayStart, workDayEnd string) (*domain.Department, error) {
    department := domain.Department{
        Name:                name,
        WorkConditionsClass: workConditionsClass,
        WorkDayStart:        workDayStart,
        WorkDayEnd:          workDayEnd,
    }
    result, err := s.repo.CreateDepartment(department)
    base.LogError(err)
    return result, err
}

func (s *service) GetDepartments() ([]domain.Department, error) {
    result, err := s.repo.GetDepartments()
    base.LogError(err)
    return result, err
}

func (s *service) UpdateDepartment(department domain.Department) (*domain.Department, error) {
    result, err := s.repo.UpdateDepartment(department)
    base.LogError(err)
    return result, err
}

func (s *service) DeleteDepartment(id int) error {
    err := s.repo.DeleteDepartment(id)
    base.LogError(err)
    return err
}

func (s *service) DepartmentById(id int) (*domain.Department, error) {
    department, err := s.repo.DepartmentById(id)
    base.LogError(err)
    return department, err
}