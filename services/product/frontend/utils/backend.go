package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nadezhda/documents/services/base"
	"github.com/nadezhda/documents/services/product/dto"
	//"github.com/levigross/grequests"
)

const (
	BACKEND_ADDRESS = "PRODUCT_BACKEND_ADDRESS"
)

// Получение адреса бэкенда
func GetBackendAddress() string {
	// получение адреса из переменной окружения или использование http://localhost:7071
	// Функция base.GetEnv принимает два аргумента (ключ и значение по умолчанию)
	// и берет значение из переменной окружения 
	// или же возвращает значение по умолчанию
	// Напишите реализацию, которая возвращает значению по умолчанию (второй аргумент)       
	address := base.GetEnv(BACKEND_ADDRESS, "http://localhost:7071")
	return address
}

// Получение единиц измерения с бэкенда
func GetUnitsFromBackend() (*dto.GetUnitsResponse, error) {
    url := fmt.Sprintf("%s/units", GetBackendAddress())
    
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    // Парсинг ответа
    var data dto.GetUnitsResponse
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return nil, err
    }
    return &data, nil
}