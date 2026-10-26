package repo

import "github.com/nadezhda/documents/services/work/domain"

type WorkRepo interface {
    GetWorks() ([]domain.Work, error)
    CreateWork(domain.Work) (*domain.Work, error)
    UpdateWork(domain.Work) (*domain.Work, error)
    DeleteWork(int) error
    WorkById(int) (*domain.Work, error)
}