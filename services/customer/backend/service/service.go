package service

import "github.com/nadezhda/documents/services/customer/domain"

type CustomerService interface {
	CreateCustomer(name string) (*domain.Customer, error)
	GetCustomers() ([]domain.Customer, error)
	CustomerById(id int) (*domain.Customer, error)
	UpdateCustomer(domain.Customer) (*domain.Customer, error)
	DeleteCustomer(id int) error
}
