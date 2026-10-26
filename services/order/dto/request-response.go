package dto

import "github.com/nadezhda/documents/services/order/domain"

// Заголовок наряда
type CreateOrderRequest struct {
    OrderNumber  string `json:"order_number"`
    OrderDate    string `json:"order_date"`
    DepartmentID int    `json:"department_id"`
    ForemanID    int    `json:"foreman_id"`
    CustomerID   int    `json:"customer_id"`
}

type CreateOrderResponse struct {
    Order *domain.Order `json:"order,omitempty"`
    Err   string        `json:"error,omitempty"`
}

type GetOrdersRequest struct{}

type GetOrdersResponse struct {
    Orders []domain.Order `json:"orders,omitempty"`
    Err    string         `json:"error,omitempty"`
}

type OrderByIdRequest struct {
    ID int `json:"id"`
}

type OrderByIdResponse struct {
    Order *domain.Order `json:"order,omitempty"`
    Err   string        `json:"error,omitempty"`
}

type UpdateOrderRequest struct {
    Order domain.Order `json:"order"`
}

type UpdateOrderResponse struct {
    Order *domain.Order `json:"order,omitempty"`
    Err   string        `json:"error,omitempty"`
}

type DeleteOrderRequest struct {
    ID int `json:"id"`
}

type DeleteOrderResponse struct {
    Err string `json:"error,omitempty"`
}

// Работы в наряде (1-я таблица)
type CreateOrderItemRequest struct {
    OrderID     int     `json:"order_id"`
    WorkID      int     `json:"work_id"`
    Description string  `json:"description"`
    Unit        string  `json:"unit"`
    Quantity    float64 `json:"quantity"`
    Rate        float64 `json:"rate"`
}

type CreateOrderItemResponse struct {
    Order *domain.Order `json:"order,omitempty"`
    Err   string        `json:"error,omitempty"`
}

type DeleteOrderItemRequest struct {
    ID int `json:"id"`
}

type DeleteOrderItemResponse struct {
    Err string `json:"error,omitempty"`
}

// Исполнители (2-я таблица)
type CreatePerformerRequest struct {
    OrderID    int     `json:"order_id"`
    EmployeeID int     `json:"employee_id"`
    KTU        float64 `json:"ktu"`
    Amount     float64 `json:"amount"`
}

type CreatePerformerResponse struct {
    Order *domain.Order `json:"order,omitempty"`
    Err   string        `json:"error,omitempty"`
}

type DeletePerformerRequest struct {
    ID int `json:"id"`
}

type DeletePerformerResponse struct {
    Err string `json:"error,omitempty"`
}