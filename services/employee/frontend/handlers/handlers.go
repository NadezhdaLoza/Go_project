package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/nadezhda/documents/services/base"
	"github.com/nadezhda/documents/services/employee/domain"
	"github.com/nadezhda/documents/services/employee/dto"
	"github.com/nadezhda/documents/services/employee/frontend/utils"
)

func EmployeesHandler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/employees", utils.GetBackendAddress())
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data dto.GetEmployeesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if data.Err != "" {
		http.Error(w, data.Err, http.StatusInternalServerError)
		return
	}
	type shortEmployee struct {
		Id       int
		Fullname string
		Post     string
	}
	employees := []shortEmployee{}
	for _, empl := range data.Employees {
		employee := shortEmployee{
			Id:       empl.Id,
			Fullname: fmt.Sprintf("%s %s %s", empl.LastName, empl.FirstName, empl.MiddleName),
			Post:     empl.Post,
		}
		employees = append(employees, employee)
	}
	employeesPage, _ := template.ParseFiles("../static/html/employees.html")
	employeesPage.Execute(w, employees)
}

func DeleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	body := dto.DeleteEmployeeRequest{Id: id}
	jsonBody, _ := json.Marshal(body)
	url := fmt.Sprintf("%s/employees", utils.GetBackendAddress())
	
	req, _ := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	var deleteResponse dto.DeleteEmployeeResponse
	if err := json.NewDecoder(resp.Body).Decode(&deleteResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if deleteResponse.Err != "" {
		http.Error(w, deleteResponse.Err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/employee/employees", http.StatusSeeOther)
}

func CreateEmployeeGetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../static/html/create-employee.html")
}

func CreateEmployeePostHandler(w http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	middleName := r.FormValue("middle_name")
	post := r.FormValue("post")
	passportSeries := r.FormValue("passport_series")
	passportNumber := r.FormValue("passport_number")
	passportIssuedBy := r.FormValue("passport_issued_by")
	passportDateOfIssue := r.FormValue("date_of_issue")
	if firstName == "" ||
		lastName == "" ||
		post == "" ||
		len(passportSeries) != 4 ||
		len(passportNumber) != 6 ||
		passportIssuedBy == "" ||
		passportDateOfIssue == "" {
		http.Error(w, base.FIELDS_VALIDATION_ERROR, http.StatusUnprocessableEntity)
		return
	}
	request := dto.CreateEmployeeRequest{
		FirstName:           firstName,
		LastName:            lastName,
		MiddleName:          middleName,
		Post:                post,
		PassportSeries:      passportSeries,
		PassportNumber:      passportNumber,
		PassportIssuedBy:    passportIssuedBy,
		PassportDateOfIssue: passportDateOfIssue,
	}
	jsonBody, _ := json.Marshal(request)
	url := fmt.Sprintf("%s/employees", utils.GetBackendAddress())
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	var data dto.CreateEmployeeResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if data.Err == "" {
		http.Redirect(w, r, "/employee/employees", http.StatusSeeOther)
	} else {
		http.Error(w, data.Err, http.StatusInternalServerError)
	}
}

func EmployeeDetailHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	url := fmt.Sprintf("%s/employees/%d", utils.GetBackendAddress(), id)
	
	requestBody := dto.EmployeeByIdRequest{Id: id}
	jsonBody, _ := json.Marshal(requestBody)
	
	req, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	var data dto.EmployeeByIdResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if data.Err != "" {
		http.Error(w, data.Err, http.StatusBadRequest)
		return
	}
	employeeDetailsPage, _ := template.ParseFiles("../static/html/employee-details.html")
	employeeDetailsPage.Execute(w, data.Employee)
}

func UpdateEmployeePostHandler(w http.ResponseWriter, r *http.Request) {
	id, idErr := strconv.Atoi(r.FormValue("id"))
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	middleName := r.FormValue("middle_name")
	post := r.FormValue("post")
	passportSeries := r.FormValue("passport_series")
	passportNumber := r.FormValue("passport_number")
	passportIssuedBy := r.FormValue("passport_issued_by")
	passportDateOfIssue := r.FormValue("date_of_issue")
	if firstName == "" ||
		lastName == "" ||
		post == "" ||
		len(passportSeries) != 4 ||
		len(passportNumber) != 6 ||
		passportIssuedBy == "" ||
		passportDateOfIssue == "" ||
		idErr != nil {

		http.Error(w, base.FIELDS_VALIDATION_ERROR, http.StatusUnprocessableEntity)
		return
	}
	request := dto.UpdateEmployeeRequest{
		Employee: domain.Employee{
			Id:                  id,
			FirstName:           firstName,
			LastName:            lastName,
			MiddleName:          middleName,
			Post:                post,
			PassportSeries:      passportSeries,
			PassportNumber:      passportNumber,
			PassportIssuedBy:    passportIssuedBy,
			PassportDateOfIssue: passportDateOfIssue,
		},
	}
	jsonBody, _ := json.Marshal(request)
	url := fmt.Sprintf("%s/employees", utils.GetBackendAddress())
	
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	var data dto.UpdateEmployeeResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if data.Err == "" {
		http.Redirect(w, r, "/employee/employees", http.StatusSeeOther)
	} else {
		http.Error(w, data.Err, http.StatusInternalServerError)
	}
}

func UpdateEmployeeGetHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	url := fmt.Sprintf("%s/employees/%d", utils.GetBackendAddress(), id)
	
	requestBody := dto.EmployeeByIdRequest{Id: id}
	jsonBody, _ := json.Marshal(requestBody)
	
	req, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	var data dto.EmployeeByIdResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if data.Err != "" {
		http.Error(w, data.Err, http.StatusBadRequest)
		return
	}
	updateEmployeePage, _ := template.ParseFiles("../static/html/update-employee.html")
	updateEmployeePage.Execute(w, data.Employee)
}