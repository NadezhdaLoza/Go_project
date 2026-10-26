package implementation

import (
    "github.com/nadezhda/documents/services/base"
    "github.com/nadezhda/documents/services/order/backend/repo"
    svc "github.com/nadezhda/documents/services/order/backend/service"
    "github.com/nadezhda/documents/services/order/domain"
)

type orderService struct {
    orderRepo repo.OrderRepo
}

func NewOrderService(orderRepo repo.OrderRepo) svc.OrderService {
    return &orderService{
        orderRepo: orderRepo,
    }
}

func (s *orderService) CreateOrder(orderNumber, orderDate string, departmentID, foremanID, customerID int) (*domain.Order, error) {
    header := domain.OrderHeader{
        OrderNumber:  orderNumber,
        OrderDate:    orderDate,
        DepartmentID: departmentID,
        ForemanID:    foremanID,
        CustomerID:   customerID,
        Status:       "draft",
    }
    
    order := domain.Order{
        OrderHeader: header,
    }
    
    result, err := s.orderRepo.CreateOrder(order)
    base.LogError(err)
    return result, err
}

func (s *orderService) GetOrders() ([]domain.Order, error) {
    result, err := s.orderRepo.GetOrders()
    base.LogError(err)
    return result, err
}

func (s *orderService) OrderById(id int) (*domain.Order, error) {
    order, err := s.orderRepo.OrderById(id)
    base.LogError(err)
    return order, err
}

func (s *orderService) UpdateOrder(order domain.Order) (*domain.Order, error) {
    result, err := s.orderRepo.UpdateOrder(order)
    base.LogError(err)
    return result, err
}

func (s *orderService) DeleteOrder(id int) error {
    err := s.orderRepo.DeleteOrder(id)
    base.LogError(err)
    return err
}

func (s *orderService) CreateOrderItem(orderID, workID int, description, unit string, quantity, rate float64) (*domain.Order, error) {
    item := domain.OrderItem{
        OrderID:     orderID,
        WorkID:      workID,
        Description: description,
        Unit:        unit,
        Quantity:    quantity,
        Rate:        rate,
    }
    result, err := s.orderRepo.CreateOrderItem(item)
    base.LogError(err)
    return result, err
}

func (s *orderService) DeleteOrderItem(itemID int) error {
    err := s.orderRepo.DeleteOrderItem(itemID)
    base.LogError(err)
    return err
}

func (s *orderService) CreatePerformer(orderID, employeeID int, ktu, amount float64) (*domain.Order, error) {
    performer := domain.Performer{
        OrderID:    orderID,
        EmployeeID: employeeID,
        KTU:        ktu,
        Amount:     amount,
    }
    result, err := s.orderRepo.CreatePerformer(performer)
    base.LogError(err)
    return result, err
}

func (s *orderService) DeletePerformer(performerID int) error {
    err := s.orderRepo.DeletePerformer(performerID)
    base.LogError(err)
    return err
}