package http

import (
    "context"
    "net/http"
    "strconv"
    
    kithttp "github.com/go-kit/kit/transport/http"
    "github.com/gorilla/mux"
    "github.com/nadezhda/documents/services/base"
    "github.com/nadezhda/documents/services/department/backend/transport"
    "github.com/nadezhda/documents/services/department/dto"
)

func NewService(
    svcEndpoints transport.Endpoints,
    options []kithttp.ServerOption,
) http.Handler {
    router := mux.NewRouter()
    errorEncoder := kithttp.ServerErrorEncoder(base.EncodeErrorResponse)
    options = append(options, errorEncoder)

    router.Methods("POST").Path("/departments").Handler(
        kithttp.NewServer(
            svcEndpoints.CreateDepartment,
            decodeCreateDepartmentRequest,
            base.EncodeResponse,
            options...,
        ))

    router.Methods("GET").Path("/departments").Handler(
        kithttp.NewServer(
            svcEndpoints.GetDepartments,
            decodeGetDepartmentsRequest,
            base.EncodeResponse,
            options...,
        ))

    router.Methods("PUT").Path("/departments").Handler(
        kithttp.NewServer(
            svcEndpoints.UpdateDepartment,
            decodeUpdateDepartmentRequest,
            base.EncodeResponse,
            options...,
        ))

    router.Methods("DELETE").Path("/departments").Handler(
        kithttp.NewServer(
            svcEndpoints.DeleteDepartment,
            decodeDeleteDepartmentRequest,
            base.EncodeResponse,
            options...,
        ))

    router.Methods("GET").Path("/departments/{id}").Handler(
        kithttp.NewServer(
            svcEndpoints.DepartmentById,
            decodeDepartmentByIdRequest,
            base.EncodeResponse,
            options...,
        ))

    return router
}

func decodeCreateDepartmentRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req dto.CreateDepartmentRequest
    err := base.DecodeBody(r, &req)
    return req, err
}

func decodeGetDepartmentsRequest(_ context.Context, r *http.Request) (interface{}, error) {
    return dto.GetDepartmentsRequest{}, nil
}

func decodeUpdateDepartmentRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req dto.UpdateDepartmentRequest
    err := base.DecodeBody(r, &req)
    return req, err
}

func decodeDeleteDepartmentRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req dto.DeleteDepartmentRequest
    err := base.DecodeBody(r, &req)
    return req, err
}

func decodeDepartmentByIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
    idStr := mux.Vars(r)["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return nil, err
    }
    req := dto.DepartmentByIdRequest{ID: id}
    return req, nil
}