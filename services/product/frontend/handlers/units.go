package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nadezhda/documents/services/base"
	"github.com/nadezhda/documents/services/product/domain"
	"github.com/nadezhda/documents/services/product/dto"
	"github.com/nadezhda/documents/services/product/frontend/utils"
)

// Обработчик страницы единиц измерения
func UnitsHandler(w http.ResponseWriter, r *http.Request) {
	// получение ед. измерения
	unitsRepo, err := utils.GetUnitsFromBackend()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if unitsRepo.Err != "" {
		http.Error(w, unitsRepo.Err, http.StatusInternalServerError)
		return
	}
	// создание шаблона
	unitsPage, err := template.ParseFiles("../static/html/units.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	unitsPage.Execute(w, unitsRepo.Units)
}

// Обработчик удаления единицы измерения
func DeleteUnitHandler(w http.ResponseWriter, r *http.Request) {
	// парсинг id
	id, _ := strconv.Atoi(r.FormValue("id"))
	// адрес ед.измерения
	unitsUrl := fmt.Sprintf("%s/units", utils.GetBackendAddress())

	// Создаем JSON запрос
	requestBody, _ := json.Marshal(map[string]interface{}{"id": id})
	
	// Создаем DELETE запрос
	req, err := http.NewRequest("DELETE", unitsUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	
	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	// Читаем ответ
	var deleteResponse dto.DeleteUnitResponse
	if err := json.NewDecoder(resp.Body).Decode(&deleteResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// проверка на ошибку
	if deleteResponse.Err != "" {
		http.Error(w, deleteResponse.Err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/product/units", http.StatusSeeOther)
}

// Обработчик страницы добавления единицы измерения
func CreateUnitGetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../static/html/create-unit.html")
}

// Обработчик добавления единицы измерения
func CreateUnitPostHandler(w http.ResponseWriter, r *http.Request) {
	// парсинг наименования
	unitName := r.FormValue("name")
	if unitName == "" {
		http.Error(w, base.FIELDS_VALIDATION_ERROR, http.StatusUnprocessableEntity)
		return
	}
	
	// создание запроса
	request := dto.CreateUnitRequest{
		Unit: unitName,
	}
	
	// адрес бэка
	unitsUrl := fmt.Sprintf("%s/units", utils.GetBackendAddress())
	
	// Создаем JSON запрос
	requestBody, _ := json.Marshal(request)
	
	// Отправляем POST запрос
	resp, err := http.Post(unitsUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	// Читаем ответ
	var data dto.CreateUnitResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// проверка на ошибку
	if data.Err == "" {
		http.Redirect(w, r, "/product/units", http.StatusSeeOther)
	} else {
		http.Error(w, data.Err, http.StatusInternalServerError)
	}
}

// Обработчик страницы изменения единицы измерения
func UpdateUnitGetHandler(w http.ResponseWriter, r *http.Request) {
	// парсинг id
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	
	// получение единицы измерения с бэка
	url := fmt.Sprintf("%s/units/%d", utils.GetBackendAddress(), id)
	
	// Используем стандартный http.Get
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	// Читаем ответ
	var unit dto.UnitByIdResponse
	if err := json.NewDecoder(resp.Body).Decode(&unit); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// unit не найден
	if unit.Unit == nil {
		http.Error(w, "the unit does not exist", http.StatusBadRequest)
		return
	}
	
	// создание шаблона
	updateUnitPage, err := template.ParseFiles("../static/html/update-unit.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	updateUnitPage.Execute(w, unit.Unit)
}

// Обработчик обновления единицы измерения
func UpdateUnitPostHandler(w http.ResponseWriter, r *http.Request) {
	// парсинг полей
	unitId, err := strconv.Atoi(r.FormValue("id"))
	unitName := r.FormValue("name")
	if unitName == "" || err != nil {
		http.Error(w, base.FIELDS_VALIDATION_ERROR, http.StatusUnprocessableEntity)
		return
	}
	
	unit := domain.Unit{
		Id:   unitId,
		Name: unitName,
	}
	
	// создание запроса
	request := dto.UpdateUnitRequest{
		Unit: unit,
	}
	
	// адрес бэка
	unitsUrl := fmt.Sprintf("%s/units", utils.GetBackendAddress())
	
	// Создаем JSON запрос
	requestBody, _ := json.Marshal(request)
	
	// Создаем PUT запрос
	req, err := http.NewRequest("PUT", unitsUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	
	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	// Читаем ответ
	var data dto.UpdateUnitResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if data.Err != "" {
		http.Error(w, data.Err, http.StatusInternalServerError)
		return
	}
	
	http.Redirect(w, r, "/product/units", http.StatusSeeOther)
}