package service

import "github.com/nadezhda/documents/services/work/domain"

type WorkService interface {
    CreateWork(name string, rate float64, unit string) (*domain.Work, error)
    GetWorks() ([]domain.Work, error)
    WorkById(id int) (*domain.Work, error)
    UpdateWork(domain.Work) (*domain.Work, error)
    DeleteWork(id int) error
}