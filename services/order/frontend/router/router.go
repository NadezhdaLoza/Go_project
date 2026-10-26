package router

import (
    "net/http"

    "github.com/gorilla/mux"
    "github.com/nadezhda/documents/services/order/frontend/handlers"
    "github.com/rs/cors"
)

func GetRouter() http.Handler {
    router := mux.NewRouter()

    // Основные маршруты нарядов
    router.HandleFunc("/documents/orders", handlers.OrdersHandler)
    router.HandleFunc("/documents/orders/create", handlers.CreateOrderGetHandler).Methods(http.MethodGet)
    router.HandleFunc("/documents/orders/create", handlers.CreateOrderPosthandler).Methods(http.MethodPost)
    router.HandleFunc("/documents/orders/update/{id:[0-9]+}", handlers.UpdateOrderGetHandler).Methods(http.MethodGet)
    router.HandleFunc("/documents/orders/update", handlers.UpdateOrderPostHandler).Methods(http.MethodPost)
    router.HandleFunc("/documents/orders/delete", handlers.DeleteOrderHandler).Methods(http.MethodPost)
    router.HandleFunc("/documents/orders/details/{id:[0-9]+}", handlers.OrderDetailsHandler).Methods(http.MethodGet)
    
    // Работы в наряде (1-е тело)
    router.HandleFunc("/documents/orders/{id:[0-9]+}/body/create", handlers.CreateOrderBodyGetHandler).Methods(http.MethodGet)
    router.HandleFunc("/documents/orders/body/create", handlers.CreateOrderBodyPostHandler).Methods(http.MethodPost)
    router.HandleFunc("/documents/orders/body/delete", handlers.DeleteOrderBodyHandler).Methods(http.MethodPost)
    
    // Исполнители (2-е тело)
    router.HandleFunc("/documents/orders/{id:[0-9]+}/performer/create", handlers.CreatePerformerGetHandler).Methods(http.MethodGet)
    router.HandleFunc("/documents/orders/performer/create", handlers.CreatePerformerPostHandler).Methods(http.MethodPost)
    router.HandleFunc("/documents/orders/performer/delete", handlers.DeletePerformerHandler).Methods(http.MethodPost)

    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../../../../static"))))
    router.PathPrefix("/documents/static/").Handler(http.StripPrefix("/documents/static/", http.FileServer(http.Dir("../../../documents/static"))))

    c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
    })

    handler := c.Handler(router)
    return handler
}