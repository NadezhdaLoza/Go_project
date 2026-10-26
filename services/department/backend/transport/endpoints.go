package transport

import (
    "context"
    
    "github.com/go-kit/kit/endpoint"
    "github.com/nadezhda/documents/services/department/backend/service"
    "github.com/nadezhda/documents/services/department/dto"
)

type Endpoints struct {
    CreateDepartment endpoint.Endpoint
    GetDepartments   endpoint.Endpoint
    UpdateDepartment endpoint.Endpoint
    DeleteDepartment endpoint.Endpoint
    DepartmentById   endpoint.Endpoint
}

func MakeEndpoints(s service.DepartmentService) Endpoints {
    return Endpoints{
        CreateDepartment: makeCreateDepartmentEndpoint(s),
        GetDepartments:   makeGetDepartmentsEndpoint(s),
        UpdateDepartment: makeUpdateDepartmentEndpoint(s),
        DeleteDepartment: makeDeleteDepartmentEndpoint(s),
        DepartmentById:   makeDepartmentByIdEndpoint(s),
    }
}

func makeCreateDepartmentEndpoint(s service.DepartmentService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        req := request.(dto.CreateDepartmentRequest)
        department, err := s.CreateDepartment(
            req.Name,
            req.WorkConditionsClass,
            req.WorkDayStart,
            req.WorkDayEnd,
        )
        return dto.CreateDepartmentResponse{Department: department}, err
    }
}

func makeGetDepartmentsEndpoint(s service.DepartmentService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        departments, err := s.GetDepartments()
        return dto.GetDepartmentsResponse{Departments: departments}, err
    }
}

func makeUpdateDepartmentEndpoint(s service.DepartmentService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        req := request.(dto.UpdateDepartmentRequest)
        department, err := s.UpdateDepartment(req.Department)
        return dto.UpdateDepartmentResponse{Department: department}, err
    }
}

func makeDeleteDepartmentEndpoint(s service.DepartmentService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        req := request.(dto.DeleteDepartmentRequest)
        return dto.DeleteDepartmentResponse{}, s.DeleteDepartment(req.ID)
    }
}

func makeDepartmentByIdEndpoint(s service.DepartmentService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        req := request.(dto.DepartmentByIdRequest)
        department, err := s.DepartmentById(req.ID)
        return dto.DepartmentByIdResponse{Department: department}, err
    }
}