package dto

import "github.com/nadezhda/documents/services/department/domain"

type CreateDepartmentRequest struct {
    Name                string `json:"name"`
    WorkConditionsClass string `json:"work_conditions_class"`
    WorkDayStart        string `json:"work_day_start"`
    WorkDayEnd          string `json:"work_day_end"`
}

type CreateDepartmentResponse struct {
    Department *domain.Department `json:"department,omitempty"`
    Err        string             `json:"error,omitempty"`
}

type GetDepartmentsRequest struct{}

type GetDepartmentsResponse struct {
    Departments []domain.Department `json:"departments,omitempty"`
    Err         string              `json:"error,omitempty"`
}

type DeleteDepartmentRequest struct {
    ID int `json:"id"`
}

type DeleteDepartmentResponse struct {
    Err string `json:"error,omitempty"`
}

type UpdateDepartmentRequest struct {
    Department domain.Department `json:"department"`
}

type UpdateDepartmentResponse struct {
    Department *domain.Department `json:"department,omitempty"`
    Err        string             `json:"error,omitempty"`
}

type DepartmentByIdRequest struct {
    ID int `json:"id"`
}

type DepartmentByIdResponse struct {
    Department *domain.Department `json:"department,omitempty"`
    Err        string             `json:"error"`
}