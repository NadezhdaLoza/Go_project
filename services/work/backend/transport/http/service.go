package http

import (
    "context"
    "net/http"
    "strconv"
    
    kithttp "github.com/go-kit/kit/transport/http"
    "github.com/gorilla/mux"
    "github.com/nadezhda/documents/services/base"
    "github.com/nadezhda/documents/services/work/backend/transport"
    "github.com/nadezhda/documents/services/work/dto"
)

func NewService(
    svcEndpoints transport.Endpoints,
    options []kithttp.ServerOption,
) http.Handler {
    router := mux.NewRouter()
    errorEncoder := kithttp.ServerErrorEncoder(base.EncodeErrorResponse)
    options = append(options, errorEncoder)

    router.Methods("POST").Path("/works").Handler(
        kithttp.NewServer(
            svcEndpoints.CreateWork,
            decodeCreateWorkRequest,
            base.EncodeResponse,
            options...,
        ))

    router.Methods("GET").Path("/works").Handler(
        kithttp.NewServer(
            svcEndpoints.GetWorks,
            decodeGetWorksRequest,
            base.EncodeResponse,
            options...,
        ))

    router.Methods("PUT").Path("/works").Handler(
        kithttp.NewServer(
            svcEndpoints.UpdateWork,
            decodeUpdateWorkRequest,
            base.EncodeResponse,
            options...,
        ))

    router.Methods("DELETE").Path("/works").Handler(
        kithttp.NewServer(
            svcEndpoints.DeleteWork,
            decodeDeleteWorkRequest,
            base.EncodeResponse,
            options...,
        ))

    router.Methods("GET").Path("/works/{id}").Handler(
        kithttp.NewServer(
            svcEndpoints.WorkById,
            decodeWorkByIdRequest,
            base.EncodeResponse,
            options...,
        ))

    return router
}

func decodeCreateWorkRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req dto.CreateWorkRequest
    err := base.DecodeBody(r, &req)
    return req, err
}

func decodeGetWorksRequest(_ context.Context, r *http.Request) (interface{}, error) {
    return dto.GetWorksRequest{}, nil
}

func decodeUpdateWorkRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req dto.UpdateWorkRequest
    err := base.DecodeBody(r, &req)
    return req, err
}

func decodeDeleteWorkRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req dto.DeleteWorkRequest
    err := base.DecodeBody(r, &req)
    return req, err
}

func decodeWorkByIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
    idStr := mux.Vars(r)["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return nil, err
    }
    req := dto.WorkByIdRequest{ID: id}
    return req, nil
}