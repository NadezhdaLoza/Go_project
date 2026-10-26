package repo

import "github.com/nadezhda/documents/services/order/domain"

type OrderRepo interface {
    // Заголовок наряда
    CreateOrder(domain.Order) (*domain.Order, error)
    GetOrders() ([]domain.Order, error)
    OrderById(id int) (*domain.Order, error)
    UpdateOrder(domain.Order) (*domain.Order, error)
    DeleteOrder(id int) error
    
    // Работы в наряде (1-я таблица)
    CreateOrderItem(domain.OrderItem) (*domain.Order, error)
    DeleteOrderItem(itemID int) error
    
    // Исполнители (2-я таблица)
    CreatePerformer(domain.Performer) (*domain.Order, error)
    DeletePerformer(performerID int) error
}