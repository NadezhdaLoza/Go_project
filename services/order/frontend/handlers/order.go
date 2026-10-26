package handlers

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "html/template"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/nadezhda/documents/services/base"
    "github.com/nadezhda/documents/services/order/domain"
    "github.com/nadezhda/documents/services/order/dto"
    "github.com/nadezhda/documents/services/order/frontend/utils"
    employeedomain "github.com/nadezhda/documents/services/employee/domain"
    customerdomain "github.com/nadezhda/documents/services/customer/domain"
    departmentdomain "github.com/nadezhda/documents/services/department/domain"
    workdomain "github.com/nadezhda/documents/services/work/domain"
)

func OrdersHandler(w http.ResponseWriter, r *http.Request) {
    url := fmt.Sprintf("%s/orders", utils.GetOrdersAddress())
    resp, err := http.Get(url)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    
    var orders dto.GetOrdersResponse
    if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    if orders.Err != "" {
        http.Error(w, orders.Err, http.StatusInternalServerError)
        return
    }

    type orderTemplate struct {
        Id            int
        OrderNumber   string
        OrderDate     string
        Organization  string
        Department    string
        Foreman       string
        TotalAmount   float64
    }

    templateData := []orderTemplate{}

    for _, order := range orders.Orders {
        
		orderDate := order.OrderHeader.OrderDate
        if len(orderDate) > 10 {
            orderDate = orderDate[:10] // "2025-12-04T00:00:00Z" -> "2025-12-04"
        }
		// Получаем организацию
        organization, _ := utils.GetCustomerById(order.OrderHeader.CustomerID)
        var organizationName string
        if organization.Err != "" {
            organizationName = "Ошибка загрузки"
        } else {
            organizationName = organization.Customer.Name
        }

        // Получаем подразделение
        department, _ := utils.GetDepartmentById(order.OrderHeader.DepartmentID)
        var departmentName string
        if department.Err != "" {
            departmentName = "Ошибка загрузки"
        } else {
            departmentName = department.Department.Name
        }

        // Получаем бригадира
        foreman, _ := utils.GetEmployeeById(order.OrderHeader.ForemanID)
        var foremanName string
        if foreman.Err != "" {
            foremanName = "Ошибка загрузки"
        } else {
            foremanName = utils.Fullname(
                foreman.Employee.LastName,
                foreman.Employee.FirstName,
                foreman.Employee.MiddleName,
            )
        }

        orderTmpl := orderTemplate{
            Id:           order.OrderHeader.ID,
            OrderNumber:  order.OrderHeader.OrderNumber,
            OrderDate:    orderDate,
            Organization: organizationName,
            Department:   departmentName,
            Foreman:      foremanName,
            TotalAmount:  order.OrderHeader.TotalAmount,
        }
        templateData = append(templateData, orderTmpl)
    }
    
    tmpl, err := template.ParseFiles("../static/html/orders.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, templateData)
}

func CreateOrderGetHandler(w http.ResponseWriter, r *http.Request) {
    type orderTemplate struct {
        Organizations []customerdomain.Customer
        Departments   []departmentdomain.Department
        Employees     []employeedomain.Employee
    }

    templateData := orderTemplate{}

    // Получаем организации
    customers, _ := utils.GetCustomers()
    if customers.Err != "" {
        base.LogError(errors.New(customers.Err))
    } else {
        templateData.Organizations = customers.Customers
    }

    // Получаем подразделения
    departments, _ := utils.GetDepartments()
    if departments.Err != "" {
        base.LogError(errors.New(departments.Err))
    } else {
        templateData.Departments = departments.Departments
    }

    // Получаем сотрудников
    employees, _ := utils.GetEmployees()
    if employees.Err != "" {
        base.LogError(errors.New(employees.Err))
    } else {
        templateData.Employees = employees.Employees
    }

    tmpl, err := template.ParseFiles("../static/html/create-order.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, templateData)
}

func CreateOrderPosthandler(w http.ResponseWriter, r *http.Request) {
    orderNumber := r.FormValue("order_number")
    orderDate := r.FormValue("order_date")
    departmentID, _ := strconv.Atoi(r.FormValue("department_id"))
    foremanID, _ := strconv.Atoi(r.FormValue("foreman_id"))
    customerID, _ := strconv.Atoi(r.FormValue("customer_id"))

    if orderNumber == "" || orderDate == "" || departmentID == 0 || foremanID == 0 || customerID == 0 {
        http.Error(w, "Все поля обязательны для заполнения", http.StatusBadRequest)
        return
    }

    order := dto.CreateOrderRequest{
        OrderNumber:  orderNumber,
        OrderDate:    orderDate,
        DepartmentID: departmentID,
        ForemanID:    foremanID,
        CustomerID:   customerID,
    }

    orderUrl := fmt.Sprintf("%s/orders", utils.GetOrdersAddress())
    
    jsonData, _ := json.Marshal(order)
    resp, err := http.Post(orderUrl, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    
    var orderResp dto.CreateOrderResponse
    if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if orderResp.Err != "" {
        http.Error(w, orderResp.Err, http.StatusInternalServerError)
        return
    }

    redirectAddress := fmt.Sprintf("/documents/orders/update/%d", orderResp.Order.OrderHeader.ID)
    http.Redirect(w, r, redirectAddress, http.StatusSeeOther)
}

func UpdateOrderGetHandler(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(mux.Vars(r)["id"])
    order, _ := utils.GetOrderById(id)
    if order.Err != "" {
        http.Error(w, order.Err, http.StatusInternalServerError)
        return
    }

    type orderItemTemplate struct {
        ID          int
        WorkID      int
        Description string
        Unit        string
        Quantity    float64
        Rate        float64
        Amount      float64
    }

    type performerTemplate struct {
        ID           int
        EmployeeID   int
        EmployeeName string
        KTU          float64
        Amount       float64
    }

    type updateOrderTemplate struct {
        Order           domain.Order
        Organizations   []customerdomain.Customer
        Departments     []departmentdomain.Department
        Employees       []employeedomain.Employee
        OrderItems      []orderItemTemplate
        Performers      []performerTemplate
        Works           []workdomain.Work
        TotalItems      float64  // Добавлено
        TotalPerformers float64  // Добавлено
    }

    templateData := updateOrderTemplate{
        Order: *order.Order,
    }

    // Заполняем список работ и вычисляем суммы
    items := []orderItemTemplate{}
    var totalItems float64
    for _, item := range order.Order.OrderBodyItems {
        amount := item.Quantity * item.Rate
        items = append(items, orderItemTemplate{
            ID:          item.ID,
            WorkID:      item.WorkID,
            Description: item.Description,
            Unit:        item.Unit,
            Quantity:    item.Quantity,
            Rate:        item.Rate,
            Amount:      amount,
        })
        totalItems += amount
    }
    templateData.OrderItems = items
    templateData.TotalItems = totalItems

    // Заполняем список исполнителей и вычисляем суммы
    performers := []performerTemplate{}
    var totalPerformers float64
    for _, performer := range order.Order.Performers {
        emp, _ := utils.GetEmployeeById(performer.EmployeeID)
        employeeName := ""
        if emp.Err == "" {
            employeeName = utils.Fullname(
                emp.Employee.LastName,
                emp.Employee.FirstName,
                emp.Employee.MiddleName,
            )
        }
        
        performers = append(performers, performerTemplate{
            ID:           performer.ID,
            EmployeeID:   performer.EmployeeID,
            EmployeeName: employeeName,
            KTU:          performer.KTU,
            Amount:       performer.Amount,
        })
        totalPerformers += performer.Amount
    }
    templateData.Performers = performers
    templateData.TotalPerformers = totalPerformers

    // Получаем данные для выпадающих списков
    customers, _ := utils.GetCustomers()
    if customers.Err == "" {
        templateData.Organizations = customers.Customers
    }

    departments, _ := utils.GetDepartments()
    if departments.Err == "" {
        templateData.Departments = departments.Departments
    }

    employees, _ := utils.GetEmployees()
    if employees.Err == "" {
        templateData.Employees = employees.Employees
    }

    works, _ := utils.GetWorks()
    if works.Err == "" {
        templateData.Works = works.Works
    }

    tmpl, err := template.ParseFiles("../static/html/update-order.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, templateData)
}

func UpdateOrderPostHandler(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.FormValue("id"))
    orderNumber := r.FormValue("order_number")
    orderDate := r.FormValue("order_date")
    departmentID, _ := strconv.Atoi(r.FormValue("department_id"))
    foremanID, _ := strconv.Atoi(r.FormValue("foreman_id"))
    customerID, _ := strconv.Atoi(r.FormValue("customer_id"))
    status := r.FormValue("status")

    if orderNumber == "" || orderDate == "" {
        http.Error(w, "Обязательные поля не заполнены", http.StatusBadRequest)
        return
    }

    header := domain.OrderHeader{
        ID:           id,
        OrderNumber:  orderNumber,
        OrderDate:    orderDate,
        DepartmentID: departmentID,
        ForemanID:    foremanID,
        CustomerID:   customerID,
        Status:       status,
    }
    
    order := domain.Order{
        OrderHeader: header,
    }

    orderRequest := dto.UpdateOrderRequest{
        Order: order,
    }

    orderUrl := fmt.Sprintf("%s/orders", utils.GetOrdersAddress())
    
    jsonData, _ := json.Marshal(orderRequest)
    req, _ := http.NewRequest("PUT", orderUrl, bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    
    var orderResp dto.UpdateOrderResponse
    if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if orderResp.Err != "" {
        http.Error(w, orderResp.Err, http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/documents/orders", http.StatusSeeOther)
}

func DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.FormValue("id"))
    body := dto.DeleteOrderRequest{ID: id}
    url := fmt.Sprintf("%s/orders", utils.GetOrdersAddress())
    
    jsonData, _ := json.Marshal(body)
    req, _ := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    
    var deleteResponse dto.DeleteOrderResponse
    if err := json.NewDecoder(resp.Body).Decode(&deleteResponse); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    if deleteResponse.Err != "" {
        http.Error(w, deleteResponse.Err, http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/documents/orders", http.StatusSeeOther)
}

func OrderDetailsHandler(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(mux.Vars(r)["id"])
    order, _ := utils.GetOrderById(id)
    if order.Err != "" {
        http.Error(w, order.Err, http.StatusInternalServerError)
        return
    }

    // Получаем дополнительные данные
    organization, _ := utils.GetCustomerById(order.Order.OrderHeader.CustomerID)
    department, _ := utils.GetDepartmentById(order.Order.OrderHeader.DepartmentID)
    foreman, _ := utils.GetEmployeeById(order.Order.OrderHeader.ForemanID)

    type orderItemTemplate struct {
        Description string
        Unit        string
        Quantity    float64
        Rate        float64
        Amount      float64
    }

    type performerTemplate struct {
        EmployeeName string
        KTU          float64
        Amount       float64
    }

    type orderDetailTemplate struct {
        OrderID         int
        OrderNumber     string
        OrderDate       string
        Organization    string
        Department      string
        Foreman         string
        Status          string
        TotalAmount     float64
        OrderItems      []orderItemTemplate
        Performers      []performerTemplate
        TotalItems      float64
        TotalPerformers float64
    }

    // Заполняем работы
    items := []orderItemTemplate{}
    var totalItems float64
    for _, item := range order.Order.OrderBodyItems {
        items = append(items, orderItemTemplate{
            Description: item.Description,
            Unit:        item.Unit,
            Quantity:    item.Quantity,
            Rate:        item.Rate,
            Amount:      item.Amount,
        })
        totalItems += item.Amount
    }

    // Заполняем исполнителей
    performers := []performerTemplate{}
    var totalPerformers float64
    for _, performer := range order.Order.Performers {
        emp, _ := utils.GetEmployeeById(performer.EmployeeID)
        employeeName := ""
        if emp.Err == "" {
            employeeName = utils.Fullname(
                emp.Employee.LastName,
                emp.Employee.FirstName,
                emp.Employee.MiddleName,
            )
        }
        
        performers = append(performers, performerTemplate{
            EmployeeName: employeeName,
            KTU:          performer.KTU,
            Amount:       performer.Amount,
        })
        totalPerformers += performer.Amount
    }

	orderDate := order.Order.OrderHeader.OrderDate
    if len(orderDate) > 10 {
        orderDate = orderDate[:10] // "2025-12-04T00:00:00Z" -> "2025-12-04"
    }
    templateData := orderDetailTemplate{
        OrderID:         order.Order.OrderHeader.ID,
        OrderNumber:     order.Order.OrderHeader.OrderNumber,
        OrderDate:       orderDate,
        Organization:    organization.Customer.Name,
        Department:      department.Department.Name,
        Foreman:         utils.Fullname(foreman.Employee.LastName, foreman.Employee.FirstName, foreman.Employee.MiddleName),
        Status:          order.Order.OrderHeader.Status,
        TotalAmount:     order.Order.OrderHeader.TotalAmount,
        OrderItems:      items,
        Performers:      performers,
        TotalItems:      totalItems,
        TotalPerformers: totalPerformers,
    }

    tmpl, err := template.ParseFiles("../static/html/order-details.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, templateData)
}