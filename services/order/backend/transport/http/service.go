package http

import (
    "context"
    "net/http"
    "strconv"
    
    kithttp "github.com/go-kit/kit/transport/http"
    "github.com/gorilla/mux"
    "github.com/nadezhda/documents/services/base"
    "github.com/nadezhda/documents/services/order/backend/transport"
    "github.com/nadezhda/documents/services/order/dto"
)

func NewService(
    svcEndpoints transport.Endpoints,
    options []kithttp.ServerOption,
) http.Handler {
    router := mux.NewRouter()
    errorEncoder := kithttp.ServerErrorEncoder(base.EncodeErrorResponse)
    options = append(options, errorEncoder)
    
    // Заголовок наряда
    router.Methods("POST").Path("/orders").Handler(
        kithttp.NewServer(
            svcEndpoints.CreateOrder,
            decodeCreateOrderRequest,
            base.EncodeResponse,
            options...,
        ))
    
    router.Methods("GET").Path("/orders").Handler(
        kithttp.NewServer(
            svcEndpoints.GetOrders,
            decodeGetOrdersRequest,
            base.EncodeResponse,
            options...,
        ))
    
    router.Methods("GET").Path("/orders/{id:[0-9]+}").Handler(
        kithttp.NewServer(
            svcEndpoints.OrderById,
            decodeOrderByIdRequest,
            base.EncodeResponse,
            options...,
        ))
    
    router.Methods("PUT").Path("/orders").Handler(
        kithttp.NewServer(
            svcEndpoints.UpdateOrder,
            decodeUpdateOrderRequest,
            base.EncodeResponse,
            options...,
        ))
    
    router.Methods("DELETE").Path("/orders").Handler(
        kithttp.NewServer(
            svcEndpoints.DeleteOrder,
            decodeDeleteOrderRequest,
            base.EncodeResponse,
            options...,
        ))
    
    // Работы в наряде (1-я таблица)
    router.Methods("POST").Path("/order-items").Handler(
        kithttp.NewServer(
            svcEndpoints.CreateOrderItem,
            decodeCreateOrderItemRequest,
            base.EncodeResponse,
            options...,
        ))
    
    router.Methods("DELETE").Path("/order-items").Handler(
        kithttp.NewServer(
            svcEndpoints.DeleteOrderItem,
            decodeDeleteOrderItemRequest,
            base.EncodeResponse,
            options...,
        ))
    
    // Исполнители (2-я таблица)
    router.Methods("POST").Path("/performers").Handler(
        kithttp.NewServer(
            svcEndpoints.CreatePerformer,
            decodeCreatePerformerRequest,
            base.EncodeResponse,
            options...,
        ))
    
    router.Methods("DELETE").Path("/performers").Handler(
        kithttp.NewServer(
            svcEndpoints.DeletePerformer,
            decodeDeletePerformerRequest,
            base.EncodeResponse,
            options...,
        ))
    
    return router
}

func decodeCreateOrderRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req dto.CreateOrderRequest
    err := base.DecodeBody(r, &req)
    return req, err
}

func decodeGetOrdersRequest(_ context.Context, r *http.Request) (interface{}, error) {
    return dto.GetOrdersRequest{}, nil
}

func decodeOrderByIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
    idStr := mux.Vars(r)["id"]
    id, err := strconv.Atoi(idStr)
    return dto.OrderByIdRequest{ID: id}, err
}

func decodeUpdateOrderRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req dto.UpdateOrderRequest
    err := base.DecodeBody(r, &req)
    return req, err
}

func decodeDeleteOrderRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req dto.DeleteOrderRequest
    err := base.DecodeBody(r, &req)
    return req, err
}

func decodeCreateOrderItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req dto.CreateOrderItemRequest
    err := base.DecodeBody(r, &req)
    return req, err
}

func decodeDeleteOrderItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req dto.DeleteOrderItemRequest
    err := base.DecodeBody(r, &req)
    return req, err
}

func decodeCreatePerformerRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req dto.CreatePerformerRequest
    err := base.DecodeBody(r, &req)
    return req, err
}

func decodeDeletePerformerRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req dto.DeletePerformerRequest
    err := base.DecodeBody(r, &req)
    return req, err
}