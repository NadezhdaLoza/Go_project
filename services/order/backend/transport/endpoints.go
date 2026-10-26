package transport

import (
    "context"
    
    "github.com/go-kit/kit/endpoint"
    "github.com/nadezhda/documents/services/order/backend/service"
    "github.com/nadezhda/documents/services/order/dto"
)

type Endpoints struct {
    // Заголовок наряда
    CreateOrder           endpoint.Endpoint
    GetOrders             endpoint.Endpoint
    OrderById             endpoint.Endpoint
    UpdateOrder           endpoint.Endpoint
    DeleteOrder           endpoint.Endpoint
    
    // Работы в наряде (1-я таблица)
    CreateOrderItem       endpoint.Endpoint
    DeleteOrderItem       endpoint.Endpoint
    
    // Исполнители (2-я таблица)
    CreatePerformer       endpoint.Endpoint
    DeletePerformer       endpoint.Endpoint
}

func MakeEndpoints(s service.OrderService) Endpoints {
    return Endpoints{
        CreateOrder:       makeCreateOrderEndpoint(s),
        GetOrders:         makeGetOrdersEndpoint(s),
        OrderById:         makeOrderByIdEndpoint(s),
        UpdateOrder:       makeUpdateOrderEndpoint(s),
        DeleteOrder:       makeDeleteOrderEndpoint(s),
        CreateOrderItem:   makeCreateOrderItemEndpoint(s),
        DeleteOrderItem:   makeDeleteOrderItemEndpoint(s),
        CreatePerformer:   makeCreatePerformerEndpoint(s),
        DeletePerformer:   makeDeletePerformerEndpoint(s),
    }
}

func makeCreateOrderEndpoint(s service.OrderService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        req := request.(dto.CreateOrderRequest)
        order, err := s.CreateOrder(
            req.OrderNumber,
            req.OrderDate,
            req.DepartmentID,
            req.ForemanID,
            req.CustomerID,
        )
        return dto.CreateOrderResponse{Order: order}, err
    }
}

func makeGetOrdersEndpoint(s service.OrderService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        orders, err := s.GetOrders()
        return dto.GetOrdersResponse{Orders: orders}, err
    }
}

func makeOrderByIdEndpoint(s service.OrderService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        req := request.(dto.OrderByIdRequest)
        order, err := s.OrderById(req.ID)
        return dto.OrderByIdResponse{Order: order}, err
    }
}

func makeUpdateOrderEndpoint(s service.OrderService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        req := request.(dto.UpdateOrderRequest)
        order, err := s.UpdateOrder(req.Order)
        return dto.UpdateOrderResponse{Order: order}, err
    }
}

func makeDeleteOrderEndpoint(s service.OrderService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        req := request.(dto.DeleteOrderRequest)
        return dto.DeleteOrderResponse{}, s.DeleteOrder(req.ID)
    }
}

func makeCreateOrderItemEndpoint(s service.OrderService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        req := request.(dto.CreateOrderItemRequest)
        order, err := s.CreateOrderItem(
            req.OrderID,
            req.WorkID,
            req.Description,
            req.Unit,
            req.Quantity,
            req.Rate,
        )
        return dto.CreateOrderItemResponse{Order: order}, err
    }
}

func makeDeleteOrderItemEndpoint(s service.OrderService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        req := request.(dto.DeleteOrderItemRequest)
        return dto.DeleteOrderItemResponse{}, s.DeleteOrderItem(req.ID)
    }
}

func makeCreatePerformerEndpoint(s service.OrderService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        req := request.(dto.CreatePerformerRequest)
        order, err := s.CreatePerformer(
            req.OrderID,
            req.EmployeeID,
            req.KTU,
            req.Amount,
        )
        return dto.CreatePerformerResponse{Order: order}, err
    }
}

func makeDeletePerformerEndpoint(s service.OrderService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        req := request.(dto.DeletePerformerRequest)
        return dto.DeletePerformerResponse{}, s.DeletePerformer(req.ID)
    }
}