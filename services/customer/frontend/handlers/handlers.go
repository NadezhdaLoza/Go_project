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
	"github.com/nadezhda/documents/services/customer/domain"
	"github.com/nadezhda/documents/services/customer/dto"
	"github.com/nadezhda/documents/services/customer/frontend/utils"
)

func CustomersHandler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/customers", utils.GetBackendAddress())
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data dto.GetCustomersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if data.Err != "" {
		http.Error(w, data.Err, http.StatusInternalServerError)
		return
	}
	
	customersPage, _ := template.ParseFiles("../static/html/customers.html")
	customersPage.Execute(w, data.Customers)
}

func DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	body := dto.DeleteCustomerRequest{Id: id}
	
	jsonBody, _ := json.Marshal(body)
	customersUrl := fmt.Sprintf("%s/customers", utils.GetBackendAddress())
	
	req, _ := http.NewRequest("DELETE", customersUrl, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	var deleteResponse dto.DeleteCustomerResponse
	if err := json.NewDecoder(resp.Body).Decode(&deleteResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if deleteResponse.Err != "" {
		http.Error(w, deleteResponse.Err, http.StatusInternalServerError)
		return
	}
	
	http.Redirect(w, r, "/customer/customers", http.StatusSeeOther)
}

func CreateCustomerGetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../static/html/create-customer.html")
}

func CreateCustomerPostHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, base.FIELDS_VALIDATION_ERROR, http.StatusUnprocessableEntity)
		return
	}
	
	request := dto.CreateCustomerRequest{
		Name: name,
	}
	
	jsonBody, _ := json.Marshal(request)
	customersUrl := fmt.Sprintf("%s/customers", utils.GetBackendAddress())
	
	resp, err := http.Post(customersUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	var data dto.CreateCustomerResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if data.Err == "" {
		http.Redirect(w, r, "/customer/customers", http.StatusSeeOther)
	} else {
		http.Error(w, data.Err, http.StatusInternalServerError)
	}
}

func UpdateCustomerGetHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	url := fmt.Sprintf("%s/customers/%d", utils.GetBackendAddress(), id)
	
	requestBody := dto.CustomerByIdRequest{Id: id}
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
	
	var data dto.CustomerByIdResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if data.Err != "" {
		http.Error(w, data.Err, http.StatusBadRequest)
		return
	}
	
	updateCustomerPage, _ := template.ParseFiles("../static/html/update-customer.html")
	updateCustomerPage.Execute(w, data.Customer)
}

func UpdateCustomerPostHandler(w http.ResponseWriter, r *http.Request) {
	customerId, _ := strconv.Atoi(r.FormValue("id"))
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, base.FIELDS_VALIDATION_ERROR, http.StatusUnprocessableEntity)
		return
	}
	
	customer := domain.Customer{
		Id:   customerId,
		Name: name,
	}
	
	request := dto.UpdateCustomerRequest{
		Customer: customer,
	}
	
	jsonBody, _ := json.Marshal(request)
	customersUrl := fmt.Sprintf("%s/customers", utils.GetBackendAddress())
	
	req, _ := http.NewRequest("PUT", customersUrl, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	var data dto.UpdateCustomerResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if data.Err == "" {
		http.Redirect(w, r, "/customer/customers", http.StatusSeeOther)
	} else {
		http.Error(w, data.Err, http.StatusInternalServerError)
	}
}