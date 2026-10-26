package repo

import "github.com/nadezhda/documents/services/department/domain"

type DepartmentRepo interface {
    GetDepartments() ([]domain.Department, error)
    CreateDepartment(domain.Department) (*domain.Department, error)
    UpdateDepartment(domain.Department) (*domain.Department, error)
    DeleteDepartment(int) error
    DepartmentById(int) (*domain.Department, error)
}