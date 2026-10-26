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
    "github.com/nadezhda/documents/services/order/dto"
    "github.com/nadezhda/documents/services/order/frontend/utils"
)

func CreateOrderBodyGetHandler(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(mux.Vars(r)["id"])
    order, _ := utils.GetOrderById(id)
    if order.Err != "" {
        http.Error(w, order.Err, http.StatusInternalServerError)
        return
    }

    works, err := utils.GetWorks()
    if err != nil {
        base.LogError(fmt.Errorf("error getting works: %v", err))
    }

    type createOrderBodyTemplate struct {
        OrderID     int
        OrderNumber string
        Works       interface{}
    }

    templateData := createOrderBodyTemplate{
        OrderID:     order.Order.OrderHeader.ID,
        OrderNumber: order.Order.OrderHeader.OrderNumber,
        Works:       works.Works,
    }

    tmpl, err := template.ParseFiles("../static/html/create-order-body.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, templateData)
}

func CreateOrderBodyPostHandler(w http.ResponseWriter, r *http.Request) {
    orderID, _ := strconv.Atoi(r.FormValue("order_id"))
    workID, _ := strconv.Atoi(r.FormValue("work_id"))
    description := r.FormValue("description")
    unit := r.FormValue("unit")
    quantity, _ := strconv.ParseFloat(r.FormValue("quantity"), 64)
    rate, _ := strconv.ParseFloat(r.FormValue("rate"), 64)

    if orderID == 0 || workID == 0 || quantity <= 0 || rate <= 0 {
        http.Error(w, "Некорректные данные", http.StatusBadRequest)
        return
    }

    data := dto.CreateOrderItemRequest{
        OrderID:     orderID,
        WorkID:      workID,
        Description: description,
        Unit:        unit,
        Quantity:    quantity,
        Rate:        rate,
    }

    url := fmt.Sprintf("%s/order-items", utils.GetOrdersAddress())
    
    jsonData, _ := json.Marshal(data)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    
    var response dto.CreateOrderItemResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    if response.Err != "" {
        http.Error(w, response.Err, http.StatusInternalServerError)
        return
    }

    redirectAddress := fmt.Sprintf("/documents/orders/update/%d", orderID)
    http.Redirect(w, r, redirectAddress, http.StatusSeeOther)
}

func DeleteOrderBodyHandler(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.FormValue("id"))
    orderID, _ := strconv.Atoi(r.FormValue("order_id"))
    
    data := dto.DeleteOrderItemRequest{
        ID: id,
    }
    
    url := fmt.Sprintf("%s/order-items", utils.GetOrdersAddress())
    
    jsonData, _ := json.Marshal(data)
    req, _ := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    
    var response dto.DeleteOrderItemResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    if response.Err != "" {
        http.Error(w, response.Err, http.StatusInternalServerError)
        return
    }

    redirectAddress := fmt.Sprintf("/documents/orders/update/%d", orderID)
    http.Redirect(w, r, redirectAddress, http.StatusSeeOther)
}