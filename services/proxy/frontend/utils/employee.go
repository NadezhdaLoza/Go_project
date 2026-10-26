package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	employeeDto "github.com/nadezhda/documents/services/employee/dto"
)

func GetEmployeeById(id int) (*employeeDto.EmployeeByIdResponse, error) {
	url := fmt.Sprintf("%s/employees/%d", GetEmployeesAddress(), id)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var employee employeeDto.EmployeeByIdResponse
	json.NewDecoder(resp.Body).Decode(&employee)
	
	return &employee, nil
}

func GetEmployees() (employeeDto.GetEmployeesResponse, error) {
	url := fmt.Sprintf("%s/employees", GetEmployeesAddress())
	
	resp, err := http.Get(url)
	if err != nil {
		return employeeDto.GetEmployeesResponse{}, err
	}
	defer resp.Body.Close()
	
	var employees employeeDto.GetEmployeesResponse
	json.NewDecoder(resp.Body).Decode(&employees)
	
	return employees, nil
}