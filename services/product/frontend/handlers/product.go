package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	//"github.com/levigross/grequests"
	"github.com/nadezhda/documents/services/product/domain"
	"github.com/nadezhda/documents/services/product/dto"
	"github.com/nadezhda/documents/services/product/frontend/utils"
)

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
    // Получение адреса сервера backend
    productsUrl := fmt.Sprintf("%s/products", utils.GetBackendAddress())
    
    // Выполнение GET-запроса на сервер backend с помощью стандартного net/http
    resp, err := http.Get(productsUrl)
    if err != nil {
        http.Error(w, "Cannot connect to backend: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    
    var products dto.GetProductsResponse
    // Декодирование ответа от сервера backend
    if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
        http.Error(w, "Invalid response from backend: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    if products.Err != "" {
        http.Error(w, products.Err, http.StatusInternalServerError)
        return
    }
    
    // Получение списка единиц измерения из backend
    units, err := utils.GetUnitsFromBackend()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if units.Err != "" {
        http.Error(w, units.Err, http.StatusInternalServerError)
        return
    }
    
    // Создание списка товаров с указанием соответствующих единиц измерения
    type ProductItemTemplate struct {
        Product domain.Product
        Unit    domain.Unit
    }
    templateItems := []ProductItemTemplate{}
    
    for _, product := range products.Products {
        templateItem := ProductItemTemplate{}
        templateItem.Product = product
        
        // Находим единицу измерения для товара
        var productUnit domain.Unit
        for _, unit := range units.Units {
            if unit.Id == product.UnitId {
                productUnit = unit
                break
            }
        }
        templateItem.Unit = productUnit
        templateItems = append(templateItems, templateItem)
    }

    // Загрузка шаблона страницы товаров
    productPage, err := template.ParseFiles("../static/html/product.html")
    if err != nil {
        http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    productPage.Execute(w, templateItems)
}


// Удаление товара
func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	productsUrl := fmt.Sprintf("%s/products", utils.GetBackendAddress())

	// Создаем JSON запрос
	requestBody, _ := json.Marshal(map[string]interface{}{"id": id})
	
	// Создаем DELETE запрос
	req, err := http.NewRequest("DELETE", productsUrl, bytes.NewBuffer(requestBody))
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
	var deleteResponse dto.DeleteProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&deleteResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if deleteResponse.Err != "" {
		http.Error(w, deleteResponse.Err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/product/products", http.StatusSeeOther)
}

// Обработчик запроса на добавление товара
func CreateProductPostHandler(w http.ResponseWriter, r *http.Request) {
	productName := r.FormValue("name")
	productPrice, priceError := strconv.Atoi(r.FormValue("price"))
	unitId, unitErr := strconv.Atoi(r.FormValue("unit_id"))
	
	if productName == "" || priceError != nil || unitErr != nil {
		http.Error(w, "fields validation error", http.StatusUnprocessableEntity)
		return
	}

	request := dto.CreateProductRequest{
		Name:   productName,
		Price:  productPrice,
		UnitId: unitId,
	}
	
	productUrl := fmt.Sprintf("%s/products", utils.GetBackendAddress())
	
	// Создаем JSON запрос
	requestBody, _ := json.Marshal(request)
	
	// Отправляем POST запрос
	resp, err := http.Post(productUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	// Читаем ответ
	var data dto.CreateProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if data.Err != "" {
		http.Error(w, data.Err, http.StatusInternalServerError)
	} else {
		http.Redirect(w, r, "/product/products", http.StatusSeeOther)
	}
}

// Обработчик страницы добавления товара
func CreateProductGetHandler(w http.ResponseWriter, r *http.Request) {
	// получить единицы измерения
	units, _ := utils.GetUnitsFromBackend()
	if units.Err != "" {
		http.Error(w, units.Err, http.StatusInternalServerError)
		return
	}
	// шаблон добавления товара
	createProductPage, _ := template.ParseFiles("../static/html/create-product.html")
	createProductPage.Execute(w, units.Units)
}

// Обработчик страницы обновления продукта
func UpdateProductGetHandler(w http.ResponseWriter, r *http.Request) {
	// парсинг id из адреса страницы
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	// получение товара с бэка
	url := fmt.Sprintf("%s/products/%d", utils.GetBackendAddress(), id)
	
	// Используем стандартный http.Get вместо grequests
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	var product dto.ProductByIdResponse
	// Декодируем JSON ответ
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// товар не найден
	if product.Err != "" {
		http.Error(w, product.Err, http.StatusBadRequest)
		return
	}
	
	// получение единиц измерения с бэка
	unitsResp, err := utils.GetUnitsFromBackend()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if unitsResp.Err != "" {
		http.Error(w, unitsResp.Err, http.StatusInternalServerError)
		return
	}
	
	// структура данных для шаблона
	type UpdateProductTemplate struct {
		Product *domain.Product
		Units   []domain.Unit
	}
	data := UpdateProductTemplate{
		Product: product.Product,
		Units:   unitsResp.Units,
	}
	
	// шаблон страницы
	updateProductPage, _ := template.ParseFiles("../static/html/update-product.html")
	updateProductPage.Execute(w, data)
}

// Обработчик запроса на обновление товара
func UpdateProductPostHandler(w http.ResponseWriter, r *http.Request) {
	productId, _ := strconv.Atoi(r.FormValue("id"))
	productName := r.FormValue("name")
	productPrice, priceError := strconv.Atoi(r.FormValue("price"))
	unitId, unitErr := strconv.Atoi(r.FormValue("unit_id"))
	
	if len(productName) == 0 || priceError != nil || unitErr != nil {
		http.Error(w, "fields_validation_error", http.StatusUnprocessableEntity)
		return
	}

	product := domain.Product{
		Id:     productId,
		Name:   productName,
		Price:  productPrice,
		UnitId: unitId,
	}

	productUrl := fmt.Sprintf("%s/products", utils.GetBackendAddress())
	request := dto.UpdateProductRequest{
		Product: product,
	}
	
	// Создаем JSON запрос
	requestBody, _ := json.Marshal(request)
	
	// Создаем PUT запрос
	req, err := http.NewRequest("PUT", productUrl, bytes.NewBuffer(requestBody))
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
	var updateResponse dto.UpdateProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if updateResponse.Err != "" {
		http.Error(w, updateResponse.Err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/product/products", http.StatusSeeOther)
}