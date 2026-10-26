package service

import "github.com/nadezhda/documents/services/order/domain"

type OrderService interface {
    // Заголовок наряда
    CreateOrder(orderNumber, orderDate string, departmentID, foremanID, customerID int) (*domain.Order, error)
    GetOrders() ([]domain.Order, error)
    OrderById(id int) (*domain.Order, error)
    UpdateOrder(domain.Order) (*domain.Order, error)
    DeleteOrder(id int) error
    
    // Работы в наряде (1-я таблица)
    CreateOrderItem(orderID, workID int, description, unit string, quantity, rate float64) (*domain.Order, error)
    DeleteOrderItem(itemID int) error
    
    // Исполнители (2-я таблица)
    CreatePerformer(orderID, employeeID int, ktu, amount float64) (*domain.Order, error)
    DeletePerformer(performerID int) error
}