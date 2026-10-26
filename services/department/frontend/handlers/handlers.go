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
	"github.com/nadezhda/documents/services/department/domain"
	"github.com/nadezhda/documents/services/department/dto"
	"github.com/nadezhda/documents/services/department/frontend/utils"
)

func DepartmentsHandler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/departments", utils.GetBackendAddress())
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data dto.GetDepartmentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if data.Err != "" {
		http.Error(w, data.Err, http.StatusInternalServerError)
		return
	}
	
	departmentsPage, _ := template.ParseFiles("../static/html/departments.html")
	departmentsPage.Execute(w, data.Departments)
}

func DeleteDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	body := dto.DeleteDepartmentRequest{ID: id}
	
	jsonBody, _ := json.Marshal(body)
	departmentsUrl := fmt.Sprintf("%s/departments", utils.GetBackendAddress())
	
	req, _ := http.NewRequest("DELETE", departmentsUrl, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	var deleteResponse dto.DeleteDepartmentResponse
	if err := json.NewDecoder(resp.Body).Decode(&deleteResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if deleteResponse.Err != "" {
		http.Error(w, deleteResponse.Err, http.StatusInternalServerError)
		return
	}
	
	http.Redirect(w, r, "/department/departments", http.StatusSeeOther)
}

func CreateDepartmentGetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../static/html/create-department.html")
}

func CreateDepartmentPostHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, base.FIELDS_VALIDATION_ERROR, http.StatusUnprocessableEntity)
		return
	}
	
	request := dto.CreateDepartmentRequest{
		Name:                name,
		WorkConditionsClass: r.FormValue("work_conditions_class"),
		WorkDayStart:        r.FormValue("work_day_start"),
		WorkDayEnd:          r.FormValue("work_day_end"),
	}
	
	jsonBody, _ := json.Marshal(request)
	departmentsUrl := fmt.Sprintf("%s/departments", utils.GetBackendAddress())
	
	resp, err := http.Post(departmentsUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	var data dto.CreateDepartmentResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if data.Err == "" {
		http.Redirect(w, r, "/department/departments", http.StatusSeeOther)
	} else {
		http.Error(w, data.Err, http.StatusInternalServerError)
	}
}

func UpdateDepartmentGetHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	url := fmt.Sprintf("%s/departments/%d", utils.GetBackendAddress(), id)
	
	requestBody := dto.DepartmentByIdRequest{ID: id}
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
	
	var data dto.DepartmentByIdResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if data.Err != "" {
		http.Error(w, data.Err, http.StatusBadRequest)
		return
	}
	
	updateDepartmentPage, _ := template.ParseFiles("../static/html/update-department.html")
	updateDepartmentPage.Execute(w, data.Department)
}

func UpdateDepartmentPostHandler(w http.ResponseWriter, r *http.Request) {
	departmentId, _ := strconv.Atoi(r.FormValue("id"))
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, base.FIELDS_VALIDATION_ERROR, http.StatusUnprocessableEntity)
		return
	}
	
	department := domain.Department{
		ID:                  departmentId,
		Name:                name,
		WorkConditionsClass: r.FormValue("work_conditions_class"),
		WorkDayStart:        r.FormValue("work_day_start"),
		WorkDayEnd:          r.FormValue("work_day_end"),
	}
	
	request := dto.UpdateDepartmentRequest{
		Department: department,
	}
	
	jsonBody, _ := json.Marshal(request)
	departmentsUrl := fmt.Sprintf("%s/departments", utils.GetBackendAddress())
	
	req, _ := http.NewRequest("PUT", departmentsUrl, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	var data dto.UpdateDepartmentResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if data.Err == "" {
		http.Redirect(w, r, "/department/departments", http.StatusSeeOther)
	} else {
		http.Error(w, data.Err, http.StatusInternalServerError)
	}
}