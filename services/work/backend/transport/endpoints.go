package transport

import (
    "context"
    
    "github.com/go-kit/kit/endpoint"
    "github.com/nadezhda/documents/services/work/backend/service"
    "github.com/nadezhda/documents/services/work/dto"
)

type Endpoints struct {
    CreateWork endpoint.Endpoint
    GetWorks   endpoint.Endpoint
    UpdateWork endpoint.Endpoint
    DeleteWork endpoint.Endpoint
    WorkById   endpoint.Endpoint
}

func MakeEndpoints(s service.WorkService) Endpoints {
    return Endpoints{
        CreateWork: makeCreateWorkEndpoint(s),
        GetWorks:   makeGetWorksEndpoint(s),
        UpdateWork: makeUpdateWorkEndpoint(s),
        DeleteWork: makeDeleteWorkEndpoint(s),
        WorkById:   makeWorkByIdEndpoint(s),
    }
}

func makeCreateWorkEndpoint(s service.WorkService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        req := request.(dto.CreateWorkRequest)
        work, err := s.CreateWork(req.Name, req.Rate, req.Unit)
        return dto.CreateWorkResponse{Work: work}, err
    }
}

func makeGetWorksEndpoint(s service.WorkService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        works, err := s.GetWorks()
        return dto.GetWorksResponse{Works: works}, err
    }
}

func makeUpdateWorkEndpoint(s service.WorkService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        req := request.(dto.UpdateWorkRequest)
        work, err := s.UpdateWork(req.Work)
        return dto.UpdateWorkResponse{Work: work}, err
    }
}

func makeDeleteWorkEndpoint(s service.WorkService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        req := request.(dto.DeleteWorkRequest)
        return dto.DeleteWorkResponse{}, s.DeleteWork(req.ID)
    }
}

func makeWorkByIdEndpoint(s service.WorkService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        req := request.(dto.WorkByIdRequest)
        work, err := s.WorkById(req.ID)
        return dto.WorkByIdResponse{Work: work}, err
    }
}