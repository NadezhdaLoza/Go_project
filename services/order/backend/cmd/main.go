package main

import (
    "fmt"
    "net/http"
    "os"

    _ "github.com/lib/pq"

    "github.com/nadezhda/documents/services/base"
    "github.com/nadezhda/documents/services/order/backend/implementation"
    "github.com/nadezhda/documents/services/order/backend/transport"
    httptransport "github.com/nadezhda/documents/services/order/backend/transport/http"

    kithttp "github.com/go-kit/kit/transport/http"
)

func main() {
    db, err := base.ConnectToDb()
    if err != nil {
        base.LogError(err)
        os.Exit(-1)
    }
    defer db.Close()

    orderRepo := implementation.NewOrderRepo(db)
    svc := implementation.NewOrderService(orderRepo)
    endpoints := transport.MakeEndpoints(svc)
    
    h := httptransport.NewService(endpoints, []kithttp.ServerOption{})
    
    fmt.Println("Listening on 7077...")
    if err := http.ListenAndServe(":7077", h); err != nil {
        base.LogError(err)
        os.Exit(-1)
    }
}