package utils

import (
    "encoding/json"
    "fmt"
    "net/http"
    
    departmentDto "github.com/nadezhda/documents/services/department/dto"
)

func GetDepartments() (*departmentDto.GetDepartmentsResponse, error) {
    url := fmt.Sprintf("%s/departments", GetDepartmentsAddress())
    
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var departments departmentDto.GetDepartmentsResponse
    if err := json.NewDecoder(resp.Body).Decode(&departments); err != nil {
        return nil, err
    }
    
    return &departments, nil
}

func GetDepartmentById(id int) (*departmentDto.DepartmentByIdResponse, error) {
    url := fmt.Sprintf("%s/departments/%d", GetDepartmentsAddress(), id)
    
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var department departmentDto.DepartmentByIdResponse
    if err := json.NewDecoder(resp.Body).Decode(&department); err != nil {
        return nil, err
    }
    
    return &department, nil
}