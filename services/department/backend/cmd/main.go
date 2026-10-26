package main

import (
    "fmt"
    "net/http"
    "os"
    
    "github.com/nadezhda/documents/services/base"
    "github.com/nadezhda/documents/services/department/backend/implementation"
    "github.com/nadezhda/documents/services/department/backend/transport"
    
    kithttp "github.com/go-kit/kit/transport/http"
    httptransport "github.com/nadezhda/documents/services/department/backend/transport/http"
)

func main() {
    db, err := base.ConnectToDb()
    if err != nil {
        base.LogError(err)
        os.Exit(-1)
    }
    
    defer db.Close()

    departmentRepo := implementation.NewDepartmentRepo(db)
    svc := implementation.NewService(departmentRepo)
    
    endpoints := transport.MakeEndpoints(svc)
    
    h := httptransport.NewService(endpoints, []kithttp.ServerOption{})
    
    fmt.Println("Listening on 7078...")
    if err := http.ListenAndServe(":7078", h); err != nil {
        base.LogError(err)
        os.Exit(-1)
    }
}