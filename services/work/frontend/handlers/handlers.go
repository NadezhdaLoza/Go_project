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
	"github.com/nadezhda/documents/services/work/domain"
	"github.com/nadezhda/documents/services/work/dto"
	"github.com/nadezhda/documents/services/work/frontend/utils"
)

func WorksHandler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/works", utils.GetBackendAddress())
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data dto.GetWorksResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if data.Err != "" {
		http.Error(w, data.Err, http.StatusInternalServerError)
		return
	}
	
	worksPage, _ := template.ParseFiles("../static/html/works.html")
	worksPage.Execute(w, data.Works)
}

func DeleteWorkHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	body := dto.DeleteWorkRequest{ID: id}
	
	jsonBody, _ := json.Marshal(body)
	worksUrl := fmt.Sprintf("%s/works", utils.GetBackendAddress())
	
	req, _ := http.NewRequest("DELETE", worksUrl, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	var deleteResponse dto.DeleteWorkResponse
	if err := json.NewDecoder(resp.Body).Decode(&deleteResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if deleteResponse.Err != "" {
		http.Error(w, deleteResponse.Err, http.StatusInternalServerError)
		return
	}
	
	http.Redirect(w, r, "/work/works", http.StatusSeeOther)
}

func CreateWorkGetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../static/html/create-work.html")
}

func CreateWorkPostHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	rateStr := r.FormValue("rate")
	if name == "" || rateStr == "" {
		http.Error(w, base.FIELDS_VALIDATION_ERROR, http.StatusUnprocessableEntity)
		return
	}
	
	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		http.Error(w, "Некорректная сумма", http.StatusUnprocessableEntity)
		return
	}
	
	request := dto.CreateWorkRequest{
		Name: name,
		Rate: rate,
		Unit: r.FormValue("unit"),
	}
	
	jsonBody, _ := json.Marshal(request)
	worksUrl := fmt.Sprintf("%s/works", utils.GetBackendAddress())
	
	resp, err := http.Post(worksUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	var data dto.CreateWorkResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if data.Err == "" {
		http.Redirect(w, r, "/work/works", http.StatusSeeOther)
	} else {
		http.Error(w, data.Err, http.StatusInternalServerError)
	}
}

func UpdateWorkGetHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	url := fmt.Sprintf("%s/works/%d", utils.GetBackendAddress(), id)
	
	requestBody := dto.WorkByIdRequest{ID: id}
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
	
	var data dto.WorkByIdResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if data.Err != "" {
		http.Error(w, data.Err, http.StatusBadRequest)
		return
	}
	
	updateWorkPage, _ := template.ParseFiles("../static/html/update-work.html")
	updateWorkPage.Execute(w, data.Work)
}

func UpdateWorkPostHandler(w http.ResponseWriter, r *http.Request) {
	workId, _ := strconv.Atoi(r.FormValue("id"))
	name := r.FormValue("name")
	rateStr := r.FormValue("rate")
	if name == "" || rateStr == "" {
		http.Error(w, base.FIELDS_VALIDATION_ERROR, http.StatusUnprocessableEntity)
		return
	}
	
	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		http.Error(w, "Некорректная сумма", http.StatusUnprocessableEntity)
		return
	}
	
	work := domain.Work{
		ID:   workId,
		Name: name,
		Rate: rate,
		Unit: r.FormValue("unit"),
	}
	
	request := dto.UpdateWorkRequest{
		Work: work,
	}
	
	jsonBody, _ := json.Marshal(request)
	worksUrl := fmt.Sprintf("%s/works", utils.GetBackendAddress())
	
	req, _ := http.NewRequest("PUT", worksUrl, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	var data dto.UpdateWorkResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if data.Err == "" {
		http.Redirect(w, r, "/work/works", http.StatusSeeOther)
	} else {
		http.Error(w, data.Err, http.StatusInternalServerError)
	}
}