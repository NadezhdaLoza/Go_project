package service

import "github.com/nadezhda/documents/services/department/domain"

type DepartmentService interface {
    CreateDepartment(name, workConditionsClass, workDayStart, workDayEnd string) (*domain.Department, error)
    GetDepartments() ([]domain.Department, error)
    DepartmentById(id int) (*domain.Department, error)
    UpdateDepartment(domain.Department) (*domain.Department, error)
    DeleteDepartment(id int) error
}