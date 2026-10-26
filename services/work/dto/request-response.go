package dto

import "github.com/nadezhda/documents/services/work/domain"

type CreateWorkRequest struct {
    Name string  `json:"name"`
    Rate float64 `json:"rate"`
    Unit string  `json:"unit"`
}

type CreateWorkResponse struct {
    Work *domain.Work `json:"work,omitempty"`
    Err  string       `json:"error,omitempty"`
}

type GetWorksRequest struct{}

type GetWorksResponse struct {
    Works []domain.Work `json:"works,omitempty"`
    Err   string        `json:"error,omitempty"`
}

type DeleteWorkRequest struct {
    ID int `json:"id"`
}

type DeleteWorkResponse struct {
    Err string `json:"error,omitempty"`
}

type UpdateWorkRequest struct {
    Work domain.Work `json:"work"`
}

type UpdateWorkResponse struct {
    Work *domain.Work `json:"work,omitempty"`
    Err  string       `json:"error,omitempty"`
}

type WorkByIdRequest struct {
    ID int `json:"id"`
}

type WorkByIdResponse struct {
    Work *domain.Work `json:"work,omitempty"`
    Err  string       `json:"error"`
}