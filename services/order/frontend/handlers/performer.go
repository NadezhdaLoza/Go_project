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

func CreatePerformerGetHandler(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(mux.Vars(r)["id"])
    order, _ := utils.GetOrderById(id)
    if order.Err != "" {
        http.Error(w, order.Err, http.StatusInternalServerError)
        return
    }

    employees, err := utils.GetEmployees()
    if err != nil {
        base.LogError(fmt.Errorf("error getting employees: %v", err))
    }

    type performerTemplate struct {
        OrderID     int
        OrderNumber string
        Employees   interface{}
    }

    templateData := performerTemplate{
        OrderID:     order.Order.OrderHeader.ID,
        OrderNumber: order.Order.OrderHeader.OrderNumber,
        Employees:   employees.Employees,
    }

    tmpl, err := template.ParseFiles("../static/html/create-performer.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, templateData)
}

func CreatePerformerPostHandler(w http.ResponseWriter, r *http.Request) {
    orderID, _ := strconv.Atoi(r.FormValue("order_id"))
    employeeID, _ := strconv.Atoi(r.FormValue("employee_id"))
    ktu, _ := strconv.ParseFloat(r.FormValue("ktu"), 64)
    amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)

    if orderID == 0 || employeeID == 0 || amount <= 0 {
        http.Error(w, "Некорректные данные", http.StatusBadRequest)
        return
    }

    data := dto.CreatePerformerRequest{
        OrderID:    orderID,
        EmployeeID: employeeID,
        KTU:        ktu,
        Amount:     amount,
    }

    url := fmt.Sprintf("%s/performers", utils.GetOrdersAddress())
    
    jsonData, _ := json.Marshal(data)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    
    var response dto.CreatePerformerResponse
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

func DeletePerformerHandler(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.FormValue("id"))
    orderID, _ := strconv.Atoi(r.FormValue("order_id"))
    
    data := dto.DeletePerformerRequest{
        ID: id,
    }
    
    url := fmt.Sprintf("%s/performers", utils.GetOrdersAddress())
    
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
    
    var response dto.DeletePerformerResponse
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